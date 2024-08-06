package meetingservice

type createInstantMeetingRes struct {
	Body struct {
		MeetingID string `json:"meeting_id" doc:"instant meeting id" example:"56ee33ddff"`
		Message   string `json:"message" example:"success"`
	}
}

type getMeetingByIdReq struct {
	MeetingID string `path:"meeting_id" minLength:"8" maxLength:"10" example:"56ee33ddff" doc:"the meeting id"`
}

type getMeetingByIdRes struct {
	Body struct {
		MeetingID string `json:"meeting_id" doc:"meeting id" example:"56ee33ddff"`
		// might return more things later
	}
}
