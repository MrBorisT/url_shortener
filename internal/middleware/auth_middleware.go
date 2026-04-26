package mw

import (
	"context"
	"net/http"
	"strings"

	"github.com/MrBorisT/url_shortener/internal/helper"
	"github.com/MrBorisT/url_shortener/internal/service"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, prefix)
			tokenString = strings.TrimSpace(tokenString)
			if tokenString == "" {
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "missing token")
				return
			}

			claims, err := authService.Verify(tokenString)
			if err != nil {
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			userID := claims.Subject
			if userID == "" {
				_ = helper.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
