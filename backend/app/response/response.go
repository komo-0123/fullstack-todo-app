package response

import (
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// エラーレスポンスをJSON形式で返す
func WriteJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": statusCode,
	})
}
