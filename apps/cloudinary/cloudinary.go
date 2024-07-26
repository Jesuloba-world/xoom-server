package cloudinary

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"

)

type Cloudinary struct {
	Cld *cloudinary.Cloudinary
}

func NewCloudinaryApp(apiKey, apiSecret, cloudName string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}
	return &Cloudinary{
		Cld: cld,
	}, nil
}
