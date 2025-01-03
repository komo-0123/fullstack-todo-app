package handler

import (
	consts "backend/app/constant"
	"backend/app/database"
	"backend/app/model"
	res "backend/app/response"
	"backend/app/validator"
	"encoding/json"
	"net/http"
)

// Todoリストをすべて取得する
func GetTodos(w http.ResponseWriter, _ *http.Request) {
	db := database.GetDB()
	rows, err := db.Query("SELECT id, title, is_complete FROM todos")
	if err != nil {
		res.WriteJsonError(w, consts.DB_ERR_FAILED_GET_TODO, http.StatusInternalServerError)
		return
	}

	var todos []model.Todo
	// レコードがある限り、次の行に進む
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete); err != nil {
			res.WriteJsonError(w, consts.DB_ERR_FAILED_GET_TODO_ROW, http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Todoリストを追加する
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		res.WriteJsonError(w, consts.INPUT_ERR_INVALID_INPUT, http.StatusBadRequest)
		return
	}

	// 入力値のバリデーション
	if err := validator.TodoInput(newTodo); err != nil {
		res.WriteJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	results, err := db.Exec("INSERT INTO todos (title, is_complete) VALUES (?, ?)", newTodo.Title, newTodo.IsComplete)
	if err != nil {
		res.WriteJsonError(w, consts.DB_ERR_FAILED_ADD_TODO, http.StatusInternalServerError)
		return
	}

	id, err := results.LastInsertId()
	if err != nil {
		res.WriteJsonError(w, consts.INPUT_ERR_FAILED_GET_ID, http.StatusInternalServerError)
		return
	}

	newTodo.Id = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}
