package logto

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bunrouter"
)

type contextKey string

const UserIdKey contextKey = "userId"

func (s *LogtoApp) AuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		huma.WriteErr(s.api, ctx, http.StatusInternalServerError, "an Error occured", fmt.Errorf("no Authorization header found"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		huma.WriteErr(s.api, ctx, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("invalid Authorization header format"))
		return
	}

	token := parts[1]

	normalContext := context.Background()
	userID, err := s.ValidateToken(normalContext, token)
	if err != nil {
		huma.WriteErr(s.api, ctx, http.StatusInternalServerError, "An error occured", err)
		return
	}

	ctx = huma.WithValue(ctx, "userId", userID)

	next(ctx)
}

func (s *LogtoApp) BunAuthMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "An error occured: no Authorization header found", http.StatusInternalServerError)
			return fmt.Errorf("no Authorization header found")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized: invalid Authorization header format", http.StatusUnauthorized)
			return fmt.Errorf("invalid Authorization header format")
		}

		token := parts[1]

		normalContext := context.Background()

		userID, err := s.ValidateToken(normalContext, token)
		if err != nil {
			http.Error(w, fmt.Sprintf("an error occured: %s", err), http.StatusInternalServerError)
			return err
		}

		ctx := context.WithValue(req.Context(), UserIdKey, userID)
		req = req.WithContext(ctx)

		return next(w, req)
	}
}
