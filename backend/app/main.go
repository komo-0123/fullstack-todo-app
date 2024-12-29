package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
    "sync"
    "io"
    "bytes"

    "golang.org/x/time/rate"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// エラーレスポンスをJSON形式で返す
func jsonError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": message,
        "status": statusCode,
    })
}
type Todo struct {
    Id int `json:"id"`
    Title string `json:"title"`
    IsComplete bool `json:"is_complete"`
}

// Todoリストをすべて取得する
func getTodos(w http.ResponseWriter, _ *http.Request) {
    rows, err := db.Query("SELECT id, title, is_complete FROM todos")
    if err != nil {
        jsonError(w, "TODOの取得に失敗しました。", http.StatusInternalServerError)
        return
    }

    var todos []Todo
    // レコードがある限り、次の行に進む
    for rows.Next() {
        var todo Todo
        if err := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete); err != nil {
            jsonError(w, "TODOの読み込みに失敗しました。", http.StatusInternalServerError)
            return
        }
        todos = append(todos, todo)
    } 

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// TodoリストのIDを指定して取得する
func getTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        jsonError(w, "IDが不正です。", http.StatusBadRequest)
    } 

    var todo Todo
    query := "SELECT id, title, is_complete FROM todos WHERE id = ?"
    err = db.QueryRow(query, id).Scan(&todo.Id, &todo.Title, &todo.IsComplete); 
    if err != nil {
        // QueryRow()は結果がない場合sql.ErrNoRowsを返すため、適切なエラーハンドリングを行う
        if err == sql.ErrNoRows {
            jsonError(w, "TODOが見つかりません。", http.StatusNotFound)
        } else {
            jsonError(w, "TODOの取得に失敗しました。", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func validationTodoInput(todo Todo) error {
    if len(strings.TrimSpace(todo.Title)) == 0 {
        return fmt.Errorf("タイトルは必須です。")
    }
    if len(todo.Title) > 255 {
        return fmt.Errorf("タイトルは255文字以内で入力してください。")
    }
    return nil
}

// Todoリストを追加する
func createTodo(w http.ResponseWriter, r *http.Request) {
    var newTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
        jsonError(w, "入力が不正です。", http.StatusBadRequest)
        return
    }

    // 入力値のバリデーション
    if err := validationTodoInput(newTodo); err != nil {
        jsonError(w, err.Error(), http.StatusBadRequest)
        return
    }

    results, err := db.Exec("INSERT INTO todos (title, is_complete) VALUES (?, ?)", newTodo.Title, newTodo.IsComplete)
    if err != nil {
        jsonError(w, "TODOの追加に失敗しました。", http.StatusInternalServerError)
        return
    }
    
    id, err := results.LastInsertId()
    if err != nil {
        jsonError(w, "IDの取得に失敗しました。", http.StatusInternalServerError)
        return
    }

    // LastInsertId()はint64型を返すので、int型に変換
    newTodo.Id = int(id)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newTodo)
}

// TodoリストのIDを指定して更新する
func updateTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        jsonError(w, "IDが不正です。", http.StatusBadRequest)
        return
    }

    var updatedTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
        jsonError(w, "入力が不正です。", http.StatusBadRequest)
        return
    }

    // 入力値のバリデーション
    if err := validationTodoInput(updatedTodo); err != nil {
        jsonError(w, err.Error(), http.StatusBadRequest)
        return
    }

    query := "UPDATE todos SET title = ?, is_complete = ? WHERE id = ?"
    result, err := db.Exec(query, updatedTodo.Title, updatedTodo.IsComplete, id)
    if err != nil {
        jsonError(w, "TODOの更新に失敗しました。", http.StatusInternalServerError)
        return
    }

    // ResultインターフェースのRowsAffected()は更新された行数を返す。
    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        jsonError(w, "更新したTODOがありません。", http.StatusNotFound)
        return
    }

    updatedTodo.Id = id
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedTodo)
}

// TodoリストのIDを指定して削除する
func deleteTodoById(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        jsonError(w, "IDが不正です。", http.StatusBadRequest)
        return
    }

    result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
    if err != nil {
        jsonError(w, "TODOの削除に失敗しました。", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        jsonError(w, "指定のTODOは削除済みです。", http.StatusNotFound)
        return
    }

    // 削除が成功した場合、ステータスコード204を返す
    w.WriteHeader(http.StatusNoContent)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getTodos(w, r)
    case http.MethodPost:
        createTodo(w, r)
    default:
        jsonError(w, "許可されていないメソッドです。", http.StatusMethodNotAllowed)
    }
}

func todoByIdHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getTodoById(w, r)
    case http.MethodPut:
        updateTodoById(w, r)
    case http.MethodDelete:
        deleteTodoById(w, r)
    default:
        jsonError(w, "許可されていないメソッドです。", http.StatusMethodNotAllowed)
    }
}

// DBの初期化
func initDB() error {
    var err error
    // dsn -> ユーザー名:パスワード@tcp(ホスト名:ポート番号)/データベース名?オプション
    dsn := "todo_user:todo_password@tcp(mysql-container:3306)/todo_db"
    db, err = sql.Open("mysql", dsn) // DBとの接続を準備
    if err != nil {
        return fmt.Errorf("failed to connect to DB: %w", err)
    }

    // DBへの接続を確認
    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping DB: %w", err)
    }

    log.Println("Connected to the database successfully!")
	return nil
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // CORSヘッダーを設定
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // プリフライトリクエスト（OPTIONS）への応答    
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func limitRequestBodyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost || r.Method == http.MethodPut {
            // 1KB までのリクエストボディのみ受け付ける
            r.Body = http.MaxBytesReader(w, r.Body, 1024) // 1024 bytes = 1KB

            // リクエストボディを読み取る
            body, err := io.ReadAll(r.Body)
            if err != nil {
                if err.Error() == "http: request body too large" {
                    jsonError(w, "リクエストボディのサイズが大きすぎます。", http.StatusRequestEntityTooLarge)
                } else {
                    jsonError(w, "リクエストボディの読み取りに失敗しました。", http.StatusInternalServerError)
                }
                return
            }

            // 読み取ったボディを再利用できるように設定
            r.Body = io.NopCloser(bytes.NewReader(body))
        }

        next.ServeHTTP(w, r)
    })
}

var limitters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    limiter, ok := limitters[ip]
    if !ok {
        limiter = rate.NewLimiter(1, 3)
        limitters[ip] = limiter
    }

    return limiter
}

func rateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        limiter := getLimiter(r.RemoteAddr)

        if !limiter.Allow() {
            jsonError(w, "リクエストが多すぎます。しばらく待ってから再度お試しください。", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func middleWareChain(next http.Handler) http.Handler {
    next = corsMiddleware(next)
    next = limitRequestBodyMiddleware(next)
    next = rateLimitMiddleware(next)
    return next
}

func main() {
    if err := initDB(); err != nil {
        log.Fatalf("failed to initialize DB: %v", err)
    }
    defer db.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("/todos", todoHandler)
    mux.HandleFunc("/todos/", todoByIdHandler)

    handlerWithMiddleWares := middleWareChain(mux)
    log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleWares))
}