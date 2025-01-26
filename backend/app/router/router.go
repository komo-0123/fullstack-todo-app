package router

import (
	"backend/app/model"
	"backend/app/response"
	"net/http"
)

func MethodRouter(method map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := method[r.Method]; ok {
			h(w, r)
		} else {
			const notAllowedMethod = "許可されていないメソッドです。"
			response.WriteTodosResponse(w, []model.Todo{}, http.StatusMethodNotAllowed, notAllowedMethod)
		}
	}
}
