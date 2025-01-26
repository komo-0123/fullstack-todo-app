package handler

import (
	"backend/app/constant"
	"backend/app/database"
	"backend/app/model"
	"backend/app/response"
	"backend/app/validator"
	"encoding/json"
	"net/http"
)

// Todoリストをすべて取得する
func GetTodos(w http.ResponseWriter, _ *http.Request) {
	db := database.GetDB()
	rows, err := db.Query("SELECT id, title, is_complete FROM todos")
	if err != nil {
		response.WriteTodosResponse(w, []model.Todo{}, http.StatusInternalServerError, constant.DB_ERR_FAILED_GET_TODO)
		return
	}

	var todos []model.Todo
	// レコードがある限り、次の行に進む
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.IsComplete); err != nil {
			response.WriteTodosResponse(w, []model.Todo{}, http.StatusInternalServerError, constant.DB_ERR_FAILED_GET_TODO_ROW)
			return
		}
		todos = append(todos, todo)
	}

	response.WriteTodosResponse(w, todos, http.StatusOK, "")
}

// Todoリストを追加する
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		response.WriteTodosResponse(w, []model.Todo{}, http.StatusBadRequest, constant.INPUT_ERR_INVALID_INPUT)
		return
	}

	// 入力値のバリデーション
	if err := validator.TodoInput(newTodo); err != nil {
		response.WriteTodosResponse(w, []model.Todo{}, http.StatusBadRequest, err.Error())
		return
	}

	db := database.GetDB()
	_, err := db.Exec("INSERT INTO todos (title, is_complete) VALUES (?, ?)", newTodo.Title, newTodo.IsComplete)
	if err != nil {
		response.WriteTodosResponse(w, []model.Todo{}, http.StatusInternalServerError, constant.DB_ERR_FAILED_ADD_TODO)
		return
	}

	response.WriteTodosResponse(w, []model.Todo{}, http.StatusCreated, "")
}
