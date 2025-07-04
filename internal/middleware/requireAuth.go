package middleware

import (
	"context"
	"Complaingo/internal/utility"
	"log"
	"net/http"
	"strings"

	appErrors "Complaingo/internal/errors"
)

type ContextKey string

const (
	ContextUserID ContextKey = "user_id"
	ContextRole   ContextKey = "role"
	ContextEmail  ContextKey = "email"
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
			log.Println("user_id not found in claims ", claims)
			WriteError(w, appErrors.ErrUnauthorized.New("Invalid user id"))
			return
		}

		email := claims["email"]
		role := claims["role"]

		ctx := context.WithValue(r.Context(), ContextUserID, int(userIdFloat))
		ctx = context.WithValue(ctx, ContextEmail, email)
		ctx = context.WithValue(ctx, ContextRole, role)

		// log.Printf("token claims: %v\n", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserId(ctx context.Context) int {
	id, ok := ctx.Value(ContextUserID).(int)
	if !ok {
		log.Println("contextUserID not found or wrong type in context")
		return 0
	}

	return id
}

func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(ContextRole).(string); ok {
		return role
	}

	return ""
}

func IsAdmin(ctx context.Context) bool {
	return GetUserRole(ctx) == "admin"
}
