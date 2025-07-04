package middleware

import (
	"log"
	"net/http"
	"strings"

	appErrors "Complaingo/internal/errors"
)

func RBAC(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r.Context())
			log.Printf("🔐 RBAC Middleware: role from context = '%s'", role)

			for _, allowed := range allowedRoles {
				if strings.EqualFold(allowed, role) {
					next.ServeHTTP(w, r)
					return
				}
			}
			WriteError(w, appErrors.ErrUnauthorized.New("Omly allowed roles can perform this action"))
		})
	}
}
