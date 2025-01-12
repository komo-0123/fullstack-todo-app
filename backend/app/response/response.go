package response

import (
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// エラーレスポンスをJSON形式で返す
func WriteJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	r := ErrorResponse{
		Error:  message,
		Status: statusCode,
	}

	json.NewEncoder(w).Encode(r)
}
