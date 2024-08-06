package meetingservice

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"

	activemeetings "github.com/Jesuloba-world/xoom-server/apps/activeMeetings"
	logto "github.com/Jesuloba-world/xoom-server/apps/logtoApp"
	humagroup "github.com/Jesuloba-world/xoom-server/lib/humaGroup"

)

type MeetingService struct {
	activeMeetings *activemeetings.ActiveMeetingService
	api            *humagroup.HumaGroup
}

func NewMeetingService(api huma.API, logto *logto.LogtoApp, activeMeetings *activemeetings.ActiveMeetingService) *MeetingService {
	return &MeetingService{
		activeMeetings: activeMeetings,
		api:            humagroup.NewHumaGroup(api, "/meetings", []string{"Meetings"}, logto.AuthMiddleware),
	}
}

func (s *MeetingService) RegisterRoutes() {
	humagroup.Post(s.api, "/instant", s.createInstantMeeting, "Create instant meeting")
	humagroup.Get(s.api, "/{meeting_id}", s.getMeetingById, "Get Meeting")
}

func (s *MeetingService) createInstantMeeting(ctx context.Context, req *struct{}) (*createInstantMeetingRes, error) {
	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	meeting, err := s.activeMeetings.CreateMeeting(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &createInstantMeetingRes{}
	resp.Body.MeetingID = meeting.ID
	resp.Body.Message = "instant meeting created"
	return resp, nil
}

func (s *MeetingService) getMeetingById(ctx context.Context, req *getMeetingByIdReq) (*getMeetingByIdRes, error) {
	activeMeeting, err := s.activeMeetings.GetMeeting(ctx, req.MeetingID)
	if err != nil {
		return nil, err
	}
	// handle scheduled meetings
	resp := &getMeetingByIdRes{}
	resp.Body.MeetingID = activeMeeting.ID
	return resp, nil
}
