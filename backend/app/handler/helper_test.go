package handler_test

import (
	"backend/app/database"
	"backend/app/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// setUpMockDBは、モックDBを作成し、それを返します。
func setUpMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックDBの作成に失敗しました: %s", err)
	}
	database.SetDB(db)

	return db, mock
}

// createTodoResponseは、テスト用のTodoResponseを作成し、それを返します。
func createTodoResponse[T model.TodoTypes](t *testing.T, data T, code int, errorMessage string) model.TodoResponse[T] {
	t.Helper()

	return model.TodoResponse[T]{
		Data: data,
		Status: model.StatusInfo{
			Code:         code,
			Error:        errorMessage != "",
			ErrorMessage: errorMessage,
		},
	}
}

// checkMockExpectationsは、モックの期待値が満たされているか確認します。
func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

// checkStatusCodeは、ステータスコードが期待値と一致しているか確認します。
func checkStatusCode(t *testing.T, want, got int) {
	t.Helper()

	if want != got {
		t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", want, got)
	}
}

// decodeResponseBodyは、レスポンスボディをデコードし、それを返します。
func decodeResponseBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()

	var got T
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("レスポンスのデコードに失敗しました: %s\nレスポンスボディ: %s", err, rec.Body.String())
	}

	return got
}

// checkResponseBodyは、レスポンスボディが期待値と一致しているか確認します。
func checkResponseBody(t *testing.T, want, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("期待したレスポンス: %v, 実際のレスポンス: %v", want, got)
	}
}

// createTestRequestは、テスト用のリクエストを作成し、それを返します。
func createTestRequest(t *testing.T, method, path, body string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}
