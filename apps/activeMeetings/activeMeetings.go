package activemeetings

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Jesuloba-world/xoom-server/util"

)

type ActiveMeetingService struct {
	rdb *redis.Client
}

type ActiveMeeting struct {
	ID        string    `json:"id"`
	CreatorID string    `json:"creator_id"`
	StartTime time.Time `json:"start_time"`
	// Description string    `json:"description"`
}

func NewActiveMeetingService(redisAddr string) *ActiveMeetingService {
	return &ActiveMeetingService{
		rdb: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
	}
}

func (s *ActiveMeetingService) CreateMeeting(ctx context.Context, creatorId string) (*ActiveMeeting, error) {
	meeting := &ActiveMeeting{
		ID:        util.GenerateMeetingID(),
		CreatorID: creatorId,
		StartTime: time.Now(),
	}

	data, err := json.Marshal(meeting)
	if err != nil {
		return nil, err
	}

	err = s.rdb.Set(ctx, "meeting:"+meeting.ID, data, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return meeting, nil
}

func (s *ActiveMeetingService) GetMeeting(ctx context.Context, meetingId string) (*ActiveMeeting, error) {
	data, err := s.rdb.Get(ctx, "meeting:"+meetingId).Bytes()
	if err != nil {
		return nil, err
	}

	meeting := new(ActiveMeeting)
	err = json.Unmarshal(data, meeting)
	if err != nil {
		return nil, err
	}

	return meeting, nil
}

func (s *ActiveMeetingService) EndMeeting(ctx context.Context, meetingId string) error {
	return s.rdb.Del(ctx, "meeting:"+meetingId).Err()
}
