package logto

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *LogtoApp) UpdateUser(ctx context.Context, userId string, updateData map[string]interface{}) error {
	payload, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal update data: %w", err)
	}

	logtoUpdateUrl := fmt.Sprintf("%s/api/users/%s", s.endpoint, userId)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, logtoUpdateUrl, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	token, err := s.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update user avatar, status code: %d", resp.StatusCode)
	}

	return nil
}
