package handler

import (
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
		res.WriteJsonError(w, "IDが不正です。", http.StatusBadRequest)
	}

	var todo model.Todo
	db := database.GetDB()
	query := "SELECT id, title, is_complete FROM todos WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&todo.Id, &todo.Title, &todo.IsComplete)
	if err != nil {
		// QueryRow()は結果がない場合sql.ErrNoRowsを返すため、適切なエラーハンドリングを行う
		if err == sql.ErrNoRows {
			res.WriteJsonError(w, "TODOが見つかりません。", http.StatusNotFound)
		} else {
			res.WriteJsonError(w, "TODOの取得に失敗しました。", http.StatusInternalServerError)
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
		res.WriteJsonError(w, "IDが不正です。", http.StatusBadRequest)
		return
	}

	var updatedTodo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		res.WriteJsonError(w, "入力が不正です。", http.StatusBadRequest)
		return
	}

	// 入力値のバリデーション
	if err := validator.TodoInput(updatedTodo); err != nil {
		res.WriteJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	query := "UPDATE todos SET title = ?, is_complete = ? WHERE id = ?"
	result, err := db.Exec(query, updatedTodo.Title, updatedTodo.IsComplete, id)
	if err != nil {
		res.WriteJsonError(w, "TODOの更新に失敗しました。", http.StatusInternalServerError)
		return
	}

	// ResultインターフェースのRowsAffected()は更新された行数を返す。
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		res.WriteJsonError(w, "更新したTODOがありません。", http.StatusNotFound)
		return
	}

	updatedTodo.Id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

// TodoリストのIDを指定して削除する
func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.WriteJsonError(w, "IDが不正です。", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		res.WriteJsonError(w, "TODOの削除に失敗しました。", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		res.WriteJsonError(w, "指定のTODOは削除済みです。", http.StatusNotFound)
		return
	}

	// 削除が成功した場合、ステータスコード204を返す
	w.WriteHeader(http.StatusNoContent)
}
