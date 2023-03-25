package rest

import (
	"net/http"
	"strings"

	"auth/custom/service"
	"auth/custom/util"
)

const (
	authHeader   = "Authorization"
	bearerSchema = "Bearer "
)

type TokenParser interface {
	ParseAccessToken(data string) (*service.AccessToken, error)
}

func Authenticate(tokenParser TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := parseBearerToken(r.Header)
			if len(token) == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			accessToken, err := tokenParser.ParseAccessToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := util.WithUserID(r.Context(), accessToken.Claims.UserID)
			ctx = util.WithUserRole(ctx, accessToken.Claims.UserRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Authorize(role service.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := util.UserRoleFromCtx(r.Context())
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !userRole.IsSufficientToRole(role) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseBearerToken(h http.Header) string {
	if h == nil {
		return ""
	}
	return strings.TrimPrefix(h.Get(authHeader), bearerSchema)
}
