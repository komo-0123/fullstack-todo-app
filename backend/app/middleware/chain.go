package middleware

import (
	"net/http"
)

// ミドルウェアを連結する
func Chain(next http.Handler, rl *RateLimiter) http.Handler {
	next = CORS(next)
	next = JSONContentType(next)
	next = LimitRequestBody(next)
	next = rl.Middleware(next)
	return next
}
