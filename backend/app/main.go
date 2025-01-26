package main

import (
	"backend/app/database"
	"backend/app/handler"
	"backend/app/middleware"
	"backend/app/router"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	serverAddress = ":8080"
)

func main() {
	initDatabase()
	defer database.GetDB().Close()

	startServer()
}

// データベースの初期化
func initDatabase() {
	if err := database.Init(); err != nil {
		log.Fatalf("failed to initialize DB: %v", err)
	}
}

// サーバーの起動
func startServer() {
	mux := setupRouter()

	lateLimiter := middleware.NewRateLimiter()
	handlerWithMiddlewares := middleware.Chain(mux, lateLimiter)

	log.Printf("Server running on http://localhost%s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, handlerWithMiddlewares))
}

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", router.MethodRouter(map[string]http.HandlerFunc{
		http.MethodGet:  handler.GetTodos,
		http.MethodPost: handler.CreateTodo,
	}))

	mux.HandleFunc("/todos/", router.MethodRouter(map[string]http.HandlerFunc{
		http.MethodGet:    handler.GetTodoById,
		http.MethodPut:    handler.UpdateTodoById,
		http.MethodDelete: handler.DeleteTodoById,
	}))

	return mux
}
