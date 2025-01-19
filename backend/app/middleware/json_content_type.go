package middleware

import (
	"backend/app/model"
	"backend/app/response"
	"net/http"
)

// JSONContentTypeは、POST/PUTリクエストのContent-TypeがJSON形式であることを確認するミドルウェア。
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if r.Header.Get("Content-Type") != "application/json" {
				const m = "JSON形式のデータを送信してください。"
				response.WriteJSON(w, []model.Todo{}, http.StatusRequestEntityTooLarge, m)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
