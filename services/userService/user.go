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
	humagroup.Post(s.api, "/profileimage", s.updateUserImage, "Update profile image")
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
