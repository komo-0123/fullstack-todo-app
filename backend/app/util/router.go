package util

import (
	consts "backend/app/constant"
	res "backend/app/response"
	"net/http"
)

func MethodRouter(method map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := method[r.Method]; ok {
			h(w, r)
		} else {
			res.WriteJsonError(w, consts.HTTP_ERR_NOT_ALLOWED_METHOD, http.StatusMethodNotAllowed)
		}
	}
}
