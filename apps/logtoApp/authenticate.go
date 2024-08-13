package logto

import (
	"context"
	"errors"
	"fmt"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func (s *LogtoApp) Authenticate(ctx context.Context) (string, error) {
	userId, ok := ctx.Value(UserIdKey).(string)
	if !ok {
		return "", errors.New("user is not authenticated")
	}
	return userId, nil
}

func (s *LogtoApp) ValidateToken(ctx context.Context, token string) (string, error) {
	jwksURL := fmt.Sprintf("%s/oidc/jwks", s.endpoint)
	keySet, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	parsedToken, err := jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(keySet),
		jwt.WithIssuer(fmt.Sprintf("%s/oidc", s.endpoint)),
		jwt.WithAudience(s.apiResourceUrl),
	)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	userID, ok := parsedToken.Get("sub")
	if !ok {
		return "", fmt.Errorf("user ID not found in token")
	}

	return userID.(string), nil
}
