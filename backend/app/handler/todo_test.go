package handler_test

import (
	"backend/app/database"
	"backend/app/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTodos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックDBの作成に失敗しました: %s", err)
	}

	database.SetDB(db)

	mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
			AddRow(1, "title1", false).
			AddRow(2, "title2", true))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)

	handler.GetTodos(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", http.StatusOK, rec.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
