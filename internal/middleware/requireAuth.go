package middleware

import (
	"context"
	"crud_api/internal/utility"
	"net/http"
	"strings"

	appErrors "crud_api/internal/errors"
)

type contextKey string

const (
	ContextUserId contextKey = "user_id"
	ContextEmail  contextKey = "email"
	ContextRole   contextKey = "role"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the bearer of the req body
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			WriteError(w, appErrors.ErrUnauthorized.New("authorization header not found"))
			return
		}
		tokeStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utility.ValidateToken(tokeStr)
		if err != nil {
			WriteError(w, appErrors.ErrUnauthorized.Wrap(err, "Claim not authorized"))
			return
		}

		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			WriteError(w, appErrors.ErrUnauthorized.New("Invalid user id"))
		}

		email := claims["email"]
		role := claims["role"]

		ctx := context.WithValue(r.Context(), ContextUserId, int(userIdFloat))
		ctx = context.WithValue(ctx, ContextEmail, email)
		ctx = context.WithValue(ctx, ContextRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
