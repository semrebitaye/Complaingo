package middleware

import (
	"net/http"
	"runtime/debug"

	appErrors "Complaingo/internal/errors"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				WriteError(w, appErrors.ErrDbFailure.New("Recovered from panic: %v\nStack trace: \n%s", r, debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
