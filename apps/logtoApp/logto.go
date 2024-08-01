package logto

import (
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-resty/resty/v2"

	"github.com/Jesuloba-world/xoom-server/apps/cloudinary"
)

type LogtoApp struct {
	endpoint          string
	applicationId     string
	applicationSecret string
	accesstoken       string
	expiresAt         time.Time
	mutex             sync.Mutex
	cloudinary        *cloudinary.Cloudinary
	api               huma.API
	apiResourceUrl    string
	client            *resty.Client
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func NewLogtoApp(endpoint, applicationId, applicationSecret string, cloudinary *cloudinary.Cloudinary, api huma.API) (*LogtoApp, error) {
	app := &LogtoApp{
		endpoint:          endpoint,
		applicationId:     applicationId,
		applicationSecret: applicationSecret,
		cloudinary:        cloudinary,
		api:               api,
		apiResourceUrl:    "http://localhost:10001",
		client:            resty.New().SetHeader("Content-Type", "application/json"),
	}

	err := app.fetchAccessToken()
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *LogtoApp) GetToken() (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	if time.Now().After(app.expiresAt) {
		err := app.fetchAccessToken()
		if err != nil {
			return "", err
		}
	}

	return app.accesstoken, nil
}
