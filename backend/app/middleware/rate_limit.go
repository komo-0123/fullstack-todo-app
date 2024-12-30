package middleware

import (
	res "backend/app/response"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var limitters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, ok := limitters[ip]
	if !ok {
		limiter = rate.NewLimiter(1, 3)
		limitters[ip] = limiter
	}

	return limiter
}

// レートリミットをかけるミドルウェア
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getLimiter(r.RemoteAddr)

		if !limiter.Allow() {
			res.WriteJsonError(w, "リクエストが多すぎます。しばらく待ってから再度お試しください。", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
