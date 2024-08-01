package logto

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s *LogtoApp) UpdateUser(ctx context.Context, userId string, updateData map[string]interface{}) error {
	logtoUpdateUrl := fmt.Sprintf("%s/api/users/%s", s.endpoint, userId)

	token, err := s.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(updateData).
		Patch(logtoUpdateUrl)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		var errorMap map[string]interface{}
		err := json.Unmarshal(resp.Body(), &errorMap)
		if err != nil {
			return fmt.Errorf("failed to update user, couldn't parse error")
		}
		return fmt.Errorf("failed to update user, message: %s", errorMap["message"].(string))
	}

	return nil
}
