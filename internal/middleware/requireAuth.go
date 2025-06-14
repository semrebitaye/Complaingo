package middleware

import (
	"context"
	"crud_api/internal/utility"
	"net/http"
	"strings"
)

func Authentiction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the bearer of the req body
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no token provided", http.StatusUnauthorized)
			return
		}
		tokeStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utility.ValidateToken(tokeStr)
		if err != nil {
			http.Error(w, "claim not authorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", int(claims["user_id"].(float64)))
		ctx = context.WithValue(ctx, "email", claims["email"])
		ctx = context.WithValue(ctx, "role", claims["role"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
