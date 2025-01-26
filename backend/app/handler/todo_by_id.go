package handler

import (
	"backend/app/constant"
	"backend/app/database"
	"backend/app/model"
	"backend/app/response"
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
		response.WriteTodoResponse(w, nil, http.StatusBadRequest, constant.INPUT_ERR_INVALID_ID)
		return
	}

	todo := &model.Todo{}
	db := database.GetDB()
	query := "SELECT id, title, is_complete FROM todos WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.IsComplete)
	if err != nil {
		// QueryRow()は結果がない場合sql.ErrNoRowsを返すため、適切なエラーハンドリングを行う
		if err == sql.ErrNoRows {
			response.WriteTodoResponse(w, nil, http.StatusNotFound, constant.DB_ERR_NOT_FOUND_TODO)
		} else {
			response.WriteTodoResponse(w, nil, http.StatusInternalServerError, constant.DB_ERR_FAILED_GET_TODO)
		}
		return
	}

	response.WriteTodoResponse(w, todo, http.StatusOK, "")
}

// TodoリストのIDを指定して更新する
func UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteTodoResponse(w, nil, http.StatusBadRequest, constant.INPUT_ERR_INVALID_ID)
		return
	}

	var updatedTodo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		response.WriteTodoResponse(w, nil, http.StatusBadRequest, constant.INPUT_ERR_INVALID_INPUT)
		return
	}

	// 入力値のバリデーション
	if err := validator.TodoInput(updatedTodo); err != nil {
		response.WriteTodoResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	db := database.GetDB()
	query := "UPDATE todos SET title = ?, is_complete = ? WHERE id = ?"
	result, err := db.Exec(query, updatedTodo.Title, updatedTodo.IsComplete, id)
	if err != nil {
		response.WriteTodoResponse(w, nil, http.StatusInternalServerError, constant.DB_ERR_FAILED_UPDATE_TODO)
		return
	}

	// ResultインターフェースのRowsAffected()は更新された行数を返す。
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		response.WriteTodoResponse(w, nil, http.StatusNotFound, constant.DB_ERR_NOT_UPDATED_TODO)
		return
	}

	response.WriteTodoResponse(w, nil, http.StatusOK, "")
}

// TodoリストのIDを指定して削除する
func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteTodoResponse(w, nil, http.StatusBadRequest, constant.INPUT_ERR_INVALID_ID)
		return
	}

	db := database.GetDB()
	result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		response.WriteTodoResponse(w, nil, http.StatusInternalServerError, constant.DB_ERR_FAILED_DELETE_TODO)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		response.WriteTodoResponse(w, nil, http.StatusNotFound, constant.DB_ERR_DELETED_TODO)
		return
	}

	response.WriteTodoResponse(w, nil, http.StatusOK, "")
}
