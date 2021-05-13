package auth

import (
	"context"
	"net/http"
	"strings"
)

type userIDContextKey string

func AuthMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			if auth == "" {
				http.Error(w, "missing bearer token in Authorization header", http.StatusUnauthorized)
				return
			}

			userId, err := ParseToken(auth)
			if err != nil {
				http.Error(w, "invalid bearer token in Authorization header", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey("user_id"), &userId)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
