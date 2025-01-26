package response

import (
	"backend/app/model"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func WriteTodoResponse(w http.ResponseWriter, todo *model.Todo, code int, errMessage string) {
	data := model.TodoResponse{
		Data: todo,
		Status: model.StatusInfo{
			Code:         code,
			Error:        errMessage != "",
			ErrorMessage: errMessage,
		},
	}

	WriteJSON(w, data, code, errMessage)
}

func WriteTodosResponse(w http.ResponseWriter, todo []model.Todo, code int, errMessage string) {
	data := model.TodosResponse{
		Data: todo,
		Status: model.StatusInfo{
			Code:         code,
			Error:        errMessage != "",
			ErrorMessage: errMessage,
		},
	}

	WriteJSON(w, data, code, errMessage)
}

type Data interface {
	model.TodoResponse | model.TodosResponse
}

// レスポンスをJSON形式で返却する
func WriteJSON[T Data](w http.ResponseWriter, data T, code int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"data": null, "status": {"code": 500, "error": true, "error_message": "内部エラーが発生しました。"}}`, http.StatusInternalServerError)
	}
}
