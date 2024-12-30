package middleware

import (
	res "backend/app/response"
	"bytes"
	"io"
	"net/http"
)

// リクエストボディのサイズを制限するミドルウェア
func LimitRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// 1KB までのリクエストボディのみ受け付ける
			r.Body = http.MaxBytesReader(w, r.Body, 1024) // 1024 bytes = 1KB

			// リクエストボディを読み取る
			body, err := io.ReadAll(r.Body)
			if err != nil {
				if err.Error() == "http: request body too large" {
					res.WriteJsonError(w, "リクエストボディのサイズが大きすぎます。", http.StatusRequestEntityTooLarge)
				} else {
					res.WriteJsonError(w, "リクエストボディの読み取りに失敗しました。", http.StatusInternalServerError)
				}
				return
			}

			// 読み取ったボディを再利用できるように設定
			r.Body = io.NopCloser(bytes.NewReader(body))
		}

		next.ServeHTTP(w, r)
	})
}
