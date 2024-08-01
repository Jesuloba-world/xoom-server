package userservice

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"

	logto "github.com/Jesuloba-world/xoom-server/apps/logtoApp"
	humagroup "github.com/Jesuloba-world/xoom-server/lib/humaGroup"
)

type UserService struct {
	api   *humagroup.HumaGroup
	logto *logto.LogtoApp
}

func NewUserService(api huma.API, logto *logto.LogtoApp) *UserService {
	return &UserService{
		api:   humagroup.NewHumaGroup(api, "/user", []string{"User management"}, logto.AuthMiddleware),
		logto: logto,
	}
}

func (s *UserService) RegisterRoutes() {
	humagroup.Patch(s.api, "/profileimage", s.updateUserImage, "Update profile image")
	humagroup.Patch(s.api, "/username", s.updateUsername, "Update username")
	humagroup.Patch(s.api, "/email", s.updateEmail, "Update email")
	humagroup.Post(s.api, "/email", s.verifyEmail, "Verify email")
}

func (s *UserService) updateUserImage(ctx context.Context, req *updateUserImageReq) (*updateUserImageRes, error) {
	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	updateData := make(map[string]interface{})

	if len(req.RawBody.File["image"]) > 0 && req.RawBody.File["image"][0] != nil {
		file := req.RawBody.File["image"][0]
		url, err := s.logto.UploadAsset(ctx, file)
		if err != nil {
			return nil, err
		}
		updateData["avatar"] = url
	} else {
		updateData["avatar"] = ""
	}

	err := s.logto.UpdateUser(ctx, userID, updateData)
	if err != nil {
		return nil, err
	}

	resp := &updateUserImageRes{}
	resp.Body.Message = "update profile image successful"
	return resp, nil
}

func (s *UserService) updateUsername(ctx context.Context, req *updateUsernameReq) (*updateUsernameRes, error) {
	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	updateData := make(map[string]interface{})

	updateData["username"] = req.Body.Username

	err := s.logto.UpdateUser(ctx, userID, updateData)
	if err != nil {
		return nil, err
	}

	resp := &updateUsernameRes{}
	resp.Body.Message = "update username successful"
	return resp, nil
}

func (s *UserService) updateEmail(ctx context.Context, req *updateEmailReq) (*updateEmailRes, error) {
	err := s.logto.SendVerificationEmail(ctx, req.Body.Email)
	if err != nil {
		return nil, err
	}

	resp := &updateEmailRes{}
	resp.Body.Message = "verify your email"
	return resp, nil
}

func (s *UserService) verifyEmail(ctx context.Context, req *verifyEmailReq) (*verifyEmailRes, error) {
	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	err := s.logto.VerifyEmail(ctx, req.Body.Email, req.Body.Otp)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	updateData["primaryEmail"] = req.Body.Email

	err = s.logto.UpdateUser(ctx, userID, updateData)
	if err != nil {
		return nil, err
	}

	resp := &verifyEmailRes{}
	resp.Body.Message = "update email successful"
	return resp, nil
}
