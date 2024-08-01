package logto

import (
	"context"
	"fmt"
	"net/http"
)

func (s *LogtoApp) SendVerificationEmail(ctx context.Context, email string) error {
	logtoUrl := fmt.Sprintf("%s/api/verification-codes", s.endpoint)

	token, err := s.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(map[string]string{
			"email": email,
		}).
		Post(logtoUrl)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed to verify email, status code: %d", resp.StatusCode())
	}

	return nil
}

func (s *LogtoApp) VerifyEmail(ctx context.Context, email string, otp string) error {
	logtoUrl := fmt.Sprintf("%s/api/verification-codes/verify", s.endpoint)

	token, err := s.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(map[string]string{
			"email":            email,
			"verificationCode": otp,
		}).
		Post(logtoUrl)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("failed to verify email, status code: %d", resp.StatusCode())
	}

	return nil
}
