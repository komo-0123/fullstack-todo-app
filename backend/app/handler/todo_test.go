package handler_test

import (
	"backend/app/database"
	"backend/app/handler"
	"backend/app/model"
	res "backend/app/response"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTodos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックDBの作成に失敗しました: %s", err)
	}
	defer db.Close()
	database.SetDB(db)

	cases := map[string]struct {
		mockSetup      func()
		wantStatusCode int
		wantBody       interface{}
		expectErr      bool
	}{
		"クエリ失敗": {
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnError(fmt.Errorf("DBエラー"))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       res.ErrorResponse{Error: "TODOの取得に失敗しました。", Status: http.StatusInternalServerError},
			expectErr:      true,
		},
		"行スキャン失敗": {
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow("不正なID", "title1", false))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       res.ErrorResponse{Error: "TODOの読み込みに失敗しました。", Status: http.StatusInternalServerError},
			expectErr:      true,
		},
		"正常系": {
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow(1, "title1", false).
						AddRow(2, "title2", true))
			},
			wantStatusCode: http.StatusOK,
			wantBody: []model.Todo{
				{ID: 1, Title: "title1", IsComplete: false},
				{ID: 2, Title: "title2", IsComplete: true},
			},
			expectErr: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.mockSetup()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)

			handler.GetTodos(rec, req)

			if rec.Code != tc.wantStatusCode {
				t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", tc.wantStatusCode, rec.Code)
			}

			if tc.expectErr {
				var got res.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("レスポンスのデコードに失敗しました: %s", err)
				}

				if !reflect.DeepEqual(got, tc.wantBody) {
					t.Errorf("期待したレスポンス: %v, 実際のレスポンス: %v", tc.wantBody, got)
				}
			} else {
				var got []model.Todo
				if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
					t.Fatalf("レスポンスのデコードに失敗しました: %s", err)
				}

				if !reflect.DeepEqual(got, tc.wantBody) {
					t.Errorf("期待したレスポンス: %v, 実際のレスポンス: %v", tc.wantBody, got)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("満たされていない期待値があります: %s", err)
			}
		})
	}
}
