package middleware

import (
	"context"
	"net/http"
	"strings"

	"letter-square-api/internal/helper"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UsernameKey contextKey = "username"

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "authorization header required"})
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := helper.ParseToken(tokenStr, jwtSecret)
			if err != nil {
				helper.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid or expired token"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) int64 {
	v, _ := r.Context().Value(UserIDKey).(int64)
	return v
}
