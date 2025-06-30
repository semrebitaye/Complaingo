package gateway

import (
	"net/http"
	"sync"
	"time"
)

var (
	rateLimitMap = make(map[string]time.Time)
	rateMutex    sync.Mutex
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr //get the client ip address to track their request count

		rateMutex.Lock()
		defer rateMutex.Unlock()

		lastRequest, exists := rateLimitMap[clientIP]
		if exists && time.Since(lastRequest) < time.Second {
			http.Error(w, "too many request", http.StatusTooManyRequests)
			return
		}
		rateLimitMap[clientIP] = time.Now()
		next.ServeHTTP(w, r)
	})
}
