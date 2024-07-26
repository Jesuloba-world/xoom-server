package logto

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (s *LogtoApp) UploadAsset(ctx context.Context, file *multipart.FileHeader) (string, error) {
	uploadresult, err := s.cloudinary.Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "xoom/profileImages",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return uploadresult.SecureURL, nil
}
