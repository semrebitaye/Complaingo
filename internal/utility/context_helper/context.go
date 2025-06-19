package contexthelper

import "context"

type ContextKey string

const (
	ContextUserId = ContextKey("user_id")
	ContextRole   = ContextKey("role")
)

func GetUserId(ctx context.Context) int {
	if id, ok := ctx.Value(ContextUserId).(int); ok {
		return id
	}

	return 0
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
