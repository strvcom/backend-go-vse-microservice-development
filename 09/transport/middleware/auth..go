package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	// AuthorizationHeader is the header key for the authorization token.
	authHeader = "Authorization"
	authSchema = "Bearer "
)

type Authenticator interface {
	VerifyToken(token string) (map[string]interface{}, error)
}

func NewAuthenticate(authenticator Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderValue := r.Header.Get(authHeader)
			if authHeaderValue == "" {
				http.Error(w, "missing auth header", http.StatusUnauthorized)
				return
			}
			t, err := parseBearerToken(authHeaderValue)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			claims, err := authenticator.VerifyToken(t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if err := verifyUserRole(claims); err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", claims["userID"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseBearerToken(header string) (string, error) {
	if !strings.HasPrefix(header, authSchema) {
		return "", errors.New("invalid auth schema")
	}

	token := strings.TrimPrefix(header, authSchema)
	if token == "" {
		return "", errors.New("empty token")
	}
	return token, nil
}

func verifyUserRole(claims map[string]interface{}) error {
	role, ok := claims["role"].(string)
	if !ok {
		return errors.New("missing role claim")
	}
	if role == "authenticated" {
		return nil
	}
	return nil
}
