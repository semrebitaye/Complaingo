package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdbRateLimit *redis.Client
var ctxRateLimit = context.Background()

const (
	rateLimitPerMinute = 5 //max req per minute per ip
	rateLimitWindow    = 1 * time.Minute
)

// setup and test the redis connection
func InitRedisRateLimit() {
	rdbRateLimit = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdbRateLimit.Ping(ctxRateLimit).Result()
	if err != nil {
		log.Fatalf("Could not init Redis for RateLimit: %v", err)
	}

	log.Println("Redis initialized for rate limiting")
}

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// use ip for rate limiting key
		ip := r.RemoteAddr

		// fixed window counter key
		key := fmt.Sprintf("rate_limit:%s", ip)

		// increment the counter for this ip
		count, err := rdbRateLimit.Incr(ctxRateLimit, key).Result()
		if err != nil {
			log.Printf("Redis INCR error for rate limiting: %v. Allowing request.", err)
			next.ServeHTTP(w, r)
			return
		}

		// if it's the first request in the window, set the expiration
		if count == 1 {
			if err := rdbRateLimit.Expire(ctxRateLimit, key, rateLimitWindow).Err(); err != nil {
				log.Printf("Redis Expire error for rate limiting: %v. Allowing request", err)
			}
		}

		if count > rateLimitPerMinute {
			// get TTL to inform the client when they can retry
			ttl, err := rdbRateLimit.TTL(ctxRateLimit, key).Result()
			if err == nil && ttl > 0 {
				w.Header().Set("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
			}
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			log.Printf("Rate limit exeeded for IP: %s (Count: %d)", ip, count)
			return
		}

		next.ServeHTTP(w, r)
	})
}
