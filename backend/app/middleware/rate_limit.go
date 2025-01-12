package middleware

import (
	"backend/app/response"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

// RateLimiterのコンストラクタ
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{limiters: make(map[string]*rate.Limiter)}
}

func (rl *RateLimiter) getLimiter(remoteAddr string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// すでにリモートアドレスに対するレートリミッターが存在する場合はそれを返す
	if limiter, exists := rl.limiters[remoteAddr]; exists {
		return limiter
	}

	const (
		limit = 1 // 1秒間に1リクエスト
		burst = 3 // バースト数
	)
	limiter := rate.NewLimiter(limit, burst)
	rl.limiters[remoteAddr] = limiter
	return limiter
}

// レートリミットを適用するミドルウェア
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := rl.getLimiter(r.RemoteAddr)

		if !limiter.Allow() {
			const tooManyRequests = "リクエストが多すぎます。しばらく待ってから再度お試しください。"
			response.WriteJSONError(w, tooManyRequests, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
