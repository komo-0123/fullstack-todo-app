package main

import (
	"backend/app/database"
	"backend/app/handler"
	"backend/app/middleware"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("failed to initialize DB: %v", err)
	}
	defer database.GetDB().Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", handler.Todo)
	mux.HandleFunc("/todos/", handler.TodoById)

	handlerWithMiddlewares := middleware.Chain(mux)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddlewares))
}
