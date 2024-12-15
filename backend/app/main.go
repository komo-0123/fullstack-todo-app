package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Todo struct {
    Id int `json:"id"`
    Title string `json:"title"`
    IsComplete bool `json:"is_completed"`
}

var todos = []Todo{
    {Id: 1, Title: "Learn Go", IsComplete: false},
    {Id: 2, Title: "Build a RESTful API", IsComplete: false},
}

// Todoリストをすべて取得する
func getTodos(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// TodoリストのIDを指定して取得する
func getTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
    } 

    for _, todo := range todos {
        if todo.Id == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(todo)
            return
        }
    }

    http.Error(w, "Todo not found", http.StatusNotFound)
}

// Todoリストを追加する
func createTodo(w http.ResponseWriter, r *http.Request) {
    var newTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    newTodo.Id = len(todos) + 1
    todos = append(todos, newTodo)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(newTodo)
}

// TodoリストのIDを指定して更新する
func updateTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var updatedTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    for i, todo := range todos {
        if todo.Id == id {
            todos[i] = updatedTodo
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(updatedTodo)
            return
        }
    }
}

// TodoリストのIDを指定して削除する
func deleteTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    for i, todo := range todos {
        if todo.Id == id {
            todos = append(todos[:i], todos[i+1:]...)

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(todos)
            return
        }
    }

    http.Error(w, "Todo not found", http.StatusNotFound)
}
// 4. 次のステップ

// 	1.	セキュリティ強化:
// 	•	バリデーションやエラーハンドリングを充実させる。
// 	•	必要であれば認証や認可を追加。
// 	2.	データベースとの連携:
// 	•	データを永続化して、コンテナ再起動後もデータが保持されるようにする。
// 	3.	テスト:
// 	•	各エンドポイントをPostmanやcurlでテスト。
// 	•	Goのnet/http/httptestを使ったユニットテスト。

func main() {
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodos(w, r) 
		case http.MethodPost:
			createTodo(w, r) 
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

    http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTodoById(w, r)
		case http.MethodPut:
			updateTodoById(w, r) 
		case http.MethodDelete:
			deleteTodoById(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

    log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}