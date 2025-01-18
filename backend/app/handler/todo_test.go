package handler_test

import (
	"backend/app/database"
	"backend/app/handler"
	"backend/app/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
	}{
		"正常系": {
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow(1, "title1", false).
						AddRow(2, "title2", true))
			},
			wantStatusCode: http.StatusOK,
			wantBody: model.TodosResponse[[]model.Todo]{
				Data: []model.Todo{
					{ID: 1, Title: "title1", IsComplete: false},
					{ID: 2, Title: "title2", IsComplete: true},
				},
				Status: model.StatusInfo{
					Code:         http.StatusOK,
					Error:        false,
					ErrorMessage: "",
				},
			},
		},
		"クエリ失敗": {
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnError(fmt.Errorf("DBエラー"))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: model.TodosResponse[[]model.Todo]{
				Data: []model.Todo{},
				Status: model.StatusInfo{
					Code:         http.StatusInternalServerError,
					Error:        true,
					ErrorMessage: "TODOの取得に失敗しました。",
				},
			},
		},
		"行スキャン失敗": {
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow("不正なID", "title1", false))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: model.TodosResponse[[]model.Todo]{
				Data: []model.Todo{},
				Status: model.StatusInfo{
					Code:         http.StatusInternalServerError,
					Error:        true,
					ErrorMessage: "TODOの読み込みに失敗しました。",
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c.mockSetup()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/todos", nil)
			req.Header.Set("Content-Type", "application/json")

			handler.GetTodos(rec, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("満たされていない期待値があります: %s", err)
			}

			if rec.Code != c.wantStatusCode {
				t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", c.wantStatusCode, rec.Code)
			}

			var got model.TodosResponse[[]model.Todo]
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("レスポンスのデコードに失敗しました: %s", err)
			}

			if !reflect.DeepEqual(got, c.wantBody) {
				t.Errorf("期待したレスポンス: %v, 実際のレスポンス: %v", c.wantBody, got)
			}
		})
	}
}

func TestCreateTodo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックDBの作成に失敗しました: %s", err)
	}
	defer db.Close()
	database.SetDB(db)

	cases := map[string]struct {
		inputBody      string
		mockSetup      func()
		wantStatusCode int
		wantBody       interface{}
	}{
		"正常系": {
			inputBody: `{"title": "新しいタスク", "is_complete": false}`,
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO todos`).
					WithArgs("新しいタスク", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantStatusCode: http.StatusCreated,
			wantBody: model.TodosResponse[[]model.Todo]{
				Data: []model.Todo{},
				Status: model.StatusInfo{
					Code:         http.StatusCreated,
					Error:        false,
					ErrorMessage: "",
				},
			},
		},
		"不正な入力": {
			inputBody:      `{"title": 123, "is_complete": false}`,
			mockSetup:      func() {},
			wantStatusCode: http.StatusBadRequest,
			wantBody: model.TodosResponse[[]model.Todo]{
				Data: []model.Todo{},
				Status: model.StatusInfo{
					Code:         http.StatusBadRequest,
					Error:        true,
					ErrorMessage: "入力が不正です。",
				},
			},
		},
		"DB追加失敗": {
			inputBody: `{"title": "新しいタスク", "is_complete": false}`,
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO todos`).
					WithArgs("新しいタスク", false).
					WillReturnError(fmt.Errorf("DBエラー"))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: model.TodosResponse[[]model.Todo]{
				Data: []model.Todo{},
				Status: model.StatusInfo{
					Code:         http.StatusInternalServerError,
					Error:        true,
					ErrorMessage: "TODOの追加に失敗しました。",
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c.mockSetup()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(c.inputBody))
			req.Header.Set("Content-Type", "application/json")

			handler.CreateTodo(rec, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("満たされていない期待値があります: %s", err)
			}

			if rec.Code != c.wantStatusCode {
				t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", c.wantStatusCode, rec.Code)
			}

			var got model.TodosResponse[[]model.Todo]
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("レスポンスのデコードに失敗しました: %s", err)
			}

			if !reflect.DeepEqual(got, c.wantBody) {
				t.Errorf("期待したレスポンス: %v, 実際のレスポンス: %v", c.wantBody, got)
			}
		})
	}
}
