package logto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (app *LogtoApp) fetchAccessToken() error {
	tokenEndpoint := fmt.Sprintf("%s/oidc/token", app.endpoint)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "all")
	data.Set("resource", fmt.Sprintf("%s/api", app.endpoint))

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", app.applicationId, app.applicationSecret)))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tokenResp tokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return err
	}

	app.accesstoken = tokenResp.AccessToken
	app.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}
