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

type updateUsernameReq struct {
	Body struct {
		Username string `json:"username" doc:"new username"`
	}
}

type updateUsernameRes struct {
	Body struct {
		Message string `json:"message" example:"success"`
	}
}

type updateEmailReq struct {
	Body struct {
		Email string `json:"email" doc:"new email"`
	}
}

type updateEmailRes struct {
	Body struct {
		Message string `json:"message" example:"success"`
	}
}

type verifyEmailReq struct {
	Body struct {
		Email string `json:"email" doc:"new email"`
		Otp   string `json:"otp" doc:"otp verification code sent to email"`
	}
}

type verifyEmailRes struct {
	Body struct {
		Message string `json:"message" example:"success"`
	}
}
