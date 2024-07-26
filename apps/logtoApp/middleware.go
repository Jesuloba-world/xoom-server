package logto

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

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
	jwksURL := fmt.Sprintf("%s/oidc/jwks", s.endpoint)
	keySet, err := jwk.Fetch(normalContext, jwksURL)
	if err != nil {
		huma.WriteErr(s.api, ctx, http.StatusInternalServerError, "Failed to fetch JWKS", err)
		return
	}

	parsedToken, err := jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(keySet),
		jwt.WithIssuer(fmt.Sprintf("%s/oidc", s.endpoint)),
		jwt.WithAudience(s.apiResourceUrl),
	)
	if err != nil {
		huma.WriteErr(s.api, ctx, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	userID, ok := parsedToken.Get("sub")
	if !ok {
		huma.WriteErr(s.api, ctx, http.StatusUnauthorized, "User ID not found in token", fmt.Errorf(""))
		return
	}

	ctx = huma.WithValue(ctx, "userId", userID)

	next(ctx)
}
