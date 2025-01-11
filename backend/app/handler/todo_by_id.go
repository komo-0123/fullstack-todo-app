package handler

import (
	consts "backend/app/constant"
	"backend/app/database"
	"backend/app/model"
	res "backend/app/response"
	"backend/app/validator"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// TodoリストのIDを指定して取得する
func GetTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.WriteJSONError(w, consts.INPUT_ERR_INVALID_ID, http.StatusBadRequest)
	}

	var todo model.Todo
	db := database.GetDB()
	query := "SELECT id, title, is_complete FROM todos WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.IsComplete)
	if err != nil {
		// QueryRow()は結果がない場合sql.ErrNoRowsを返すため、適切なエラーハンドリングを行う
		if err == sql.ErrNoRows {
			res.WriteJSONError(w, consts.DB_ERR_NOT_FOUND_TODO, http.StatusNotFound)
		} else {
			res.WriteJSONError(w, consts.DB_ERR_FAILED_GET_TODO, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// TodoリストのIDを指定して更新する
func UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.WriteJSONError(w, consts.INPUT_ERR_INVALID_ID, http.StatusBadRequest)
		return
	}

	var updatedTodo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		res.WriteJSONError(w, consts.INPUT_ERR_INVALID_INPUT, http.StatusBadRequest)
		return
	}

	// 入力値のバリデーション
	if err := validator.TodoInput(updatedTodo); err != nil {
		res.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	query := "UPDATE todos SET title = ?, is_complete = ? WHERE id = ?"
	result, err := db.Exec(query, updatedTodo.Title, updatedTodo.IsComplete, id)
	if err != nil {
		res.WriteJSONError(w, consts.DB_ERR_FAILED_UPDATE_TODO, http.StatusInternalServerError)
		return
	}

	// ResultインターフェースのRowsAffected()は更新された行数を返す。
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		res.WriteJSONError(w, consts.DB_ERR_NOT_UPDATED_TODO, http.StatusNotFound)
		return
	}

	updatedTodo.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

// TodoリストのIDを指定して削除する
func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.WriteJSONError(w, consts.INPUT_ERR_INVALID_ID, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		res.WriteJSONError(w, consts.DB_ERR_FAILED_DELETE_TODO, http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		res.WriteJSONError(w, consts.DB_ERR_DELETED_TODO, http.StatusNotFound)
		return
	}

	// 削除が成功した場合、ステータスコード204を返す
	w.WriteHeader(http.StatusNoContent)
}
