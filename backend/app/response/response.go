package response

import (
	"backend/app/model"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// レスポンスをJSON形式で返却する
func WriteJSON[T model.Todo | []model.Todo](w http.ResponseWriter, data T, code int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := model.TodosResponse[T]{
		Data: data,
		Status: model.StatusInfo{
			Code:         code,
			Error:        errMessage != "",
			ErrorMessage: errMessage,
		},
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"data": null, "status": {"code": 500, "error": true, "error_message": "内部エラーが発生しました。"}}`, http.StatusInternalServerError)
	}
}
