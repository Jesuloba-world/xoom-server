package userservice

import "mime/multipart"

type updateUserImageReq struct {
	Body struct {
		Image multipart.File `json:"image" doc:"upload image file, use formdata"`
	}
	RawBody multipart.Form
}

type updateUserImageRes struct {
	Body struct {
		Message string `json:"message" example:"success"`
	}
}
