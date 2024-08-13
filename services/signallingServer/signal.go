package signallingserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bunrouter"

	activemeetings "github.com/Jesuloba-world/xoom-server/apps/activeMeetings"
	logto "github.com/Jesuloba-world/xoom-server/apps/logtoApp"
	"github.com/Jesuloba-world/xoom-server/lib/httpError"
)

type SignallingServer struct {
	upgrader       websocket.Upgrader
	rdb            *redis.Client
	activeMeetings *activemeetings.ActiveMeetingService
	logto          *logto.LogtoApp
	clients        map[string]map[string]*websocket.Conn
	clientMu       sync.RWMutex
}

func NewSignallingServer(rdb *redis.Client, activeMeetings *activemeetings.ActiveMeetingService, logto *logto.LogtoApp) *SignallingServer {
	s := &SignallingServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// we should check headers here and allow only specified domains
				return true
			},
		},
		rdb:            rdb,
		activeMeetings: activeMeetings,
		logto:          logto,
		clients:        make(map[string]map[string]*websocket.Conn),
	}
	go s.subscribeToRedis()
	return s
}

func (s *SignallingServer) handleWebsocket(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	meetingID := req.Param("meetingId")
	if meetingID == "" {
		return bunrouter.JSON(w, httpError.HTTPError{StatusCode: http.StatusBadRequest, Code: "missing_id", Message: "meetingId is required"})
	}

	slog.Info(meetingID)

	token := req.URL.Query().Get("token")
	if token == "" {
		return bunrouter.JSON(w, httpError.HTTPError{StatusCode: http.StatusUnauthorized, Code: "missing_token", Message: "Token is required"})
	}

	// only authenticated users can join meetings
	userId, err := s.logto.ValidateToken(ctx, token)
	if err != nil {
		return bunrouter.JSON(w, httpError.HTTPError{StatusCode: http.StatusUnauthorized, Code: "unauthorized", Message: "Authentication failed"})
	}

	_, err = s.activeMeetings.GetMeeting(ctx, meetingID)
	if err != nil {
		if err == redis.Nil {
			return bunrouter.JSON(w, httpError.HTTPError{StatusCode: http.StatusNotFound, Code: "not_found", Message: "Meeting not found"})
		} else {
			return bunrouter.JSON(w, httpError.HTTPError{StatusCode: http.StatusInternalServerError, Code: "internal_error", Message: "Failed to get meeting"})
		}
	}

	conn, err := s.upgrader.Upgrade(w, req.Request, nil)
	if err != nil {
		slog.Error("Error upgrading to Websocket", "error", err)
		return err
	}

	s.addClient(meetingID, userId, conn)

	go s.handleClient(ctx, meetingID, userId, conn)
	return nil
}

func (s *SignallingServer) handleClient(ctx context.Context, meetingId, userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		s.removeClient(meetingId, userID)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Error reading message", "err", err)
			break
		}

		var signal struct {
			Type string `json:"type"`
			To   string `json:"to,omitempty"`
		}
		if err := json.Unmarshal(msg, &signal); err != nil {
			slog.Error("Error parsing signal", "err", err)
			continue
		}

		switch signal.Type {
		case "offer", "answer", "ice-candidate":
			s.relaysignal(ctx, meetingId, signal.To, msg)
		case "join":
			s.broadcastNewPeer(ctx, meetingId, userID)
		}
	}
}

func (s *SignallingServer) relaysignal(ctx context.Context, meetingId, to string, msg []byte) {
	channel := "meeting:" + meetingId + ":client:" + to
	err := s.rdb.Publish(ctx, channel, msg).Err()
	if err != nil {
		slog.Error("Error relaying signal", "err", err)
	}
}

func (s *SignallingServer) broadcastNewPeer(ctx context.Context, meetingId, userId string) {
	message := map[string]string{
		"type":   "new-peer",
		"peerId": userId,
	}
	jsonMsg, _ := json.Marshal(message)
	channel := "meeting:" + meetingId + ":broadcast"
	err := s.rdb.Publish(ctx, channel, jsonMsg).Err()
	if err != nil {
		slog.Error("Error broadcasting new peer", "err", err)
	}
}

func (s *SignallingServer) addClient(meetingId, userId string, conn *websocket.Conn) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	if _, ok := s.clients[meetingId]; !ok {
		s.clients[meetingId] = make(map[string]*websocket.Conn)
	}
	s.clients[meetingId][userId] = conn
}

func (s *SignallingServer) removeClient(meetingId, userId string) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()

	ctx := context.Background()
	if meeting, ok := s.clients[meetingId]; ok {
		delete(meeting, userId)
		if len(meeting) == 0 {
			delete(s.clients, meetingId)
			s.activeMeetings.EndMeeting(ctx, meetingId)
		}
	}
}

func (s *SignallingServer) subscribeToRedis() {
	pubsub := s.rdb.PSubscribe(context.Background(), "meeting:*:client:*", "meeting:*:broadcast")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		s.handleRedisMessage(msg)
	}
}

func (s *SignallingServer) handleRedisMessage(msg *redis.Message) {
	parts := strings.Split(msg.Channel, ":")
	if len(parts) < 3 {
		slog.Error("Invalid channel format", "channel", msg.Channel)
	}

	meetingId := parts[1]
	if parts[2] == "broadcast" {
		s.broadcastToMeeting(meetingId, []byte(msg.Payload))
	} else if parts[2] == "client" && len(parts) == 4 {
		userId := parts[3]
		s.sendToClient(meetingId, userId, []byte(msg.Payload))
	}
}

func (s *SignallingServer) broadcastToMeeting(meetingId string, msg []byte) {
	s.clientMu.RLock()
	defer s.clientMu.RUnlock()

	if meeting, ok := s.clients[meetingId]; ok {
		for _, conn := range meeting {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				slog.Error("Error sending message to client", "err", err)
			}
		}
	}
}

func (s *SignallingServer) sendToClient(meetingId, userId string, msg []byte) {
	s.clientMu.RLock()
	defer s.clientMu.RUnlock()

	if meeting, ok := s.clients[meetingId]; ok {
		if conn, ok := meeting[userId]; ok {
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				slog.Error("Error sending message to client", "err", err)
			}
		}
	}
}

func (s *SignallingServer) RegisterRoute(router *bunrouter.Router) {
	router.GET("/room/:meetingId", s.handleWebsocket)
}
