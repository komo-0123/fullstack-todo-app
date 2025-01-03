package middleware

import (
	"net/http"
)

// ミドルウェアを連結する
func Chain(next http.Handler) http.Handler {
	next = CORS(next)
	next = LimitRequestBody(next)
	next = RateLimitMiddleware(next)
	return next
}
