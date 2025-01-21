package handler_test

import (
	"backend/app/database"
	"backend/app/handler"
	"backend/app/model"
	"database/sql"
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
	db, mock := setUpMockDB(t)
	defer db.Close()

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
			wantBody: model.TodoResponse[[]model.Todo]{
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
			wantBody: model.TodoResponse[[]model.Todo]{
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
			wantBody: model.TodoResponse[[]model.Todo]{
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
			req := createTestRequest(t, http.MethodGet, "/todos", "")

			handler.GetTodos(rec, req)

			checkMockExpectations(t, mock)
			checkStatusCode(t, c.wantStatusCode, rec.Code)
			got := decodeResponseBody[model.TodoResponse[[]model.Todo]](t, rec)
			checkResponseBody(t, c.wantBody, got)
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
			wantBody: model.TodoResponse[[]model.Todo]{
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
			wantBody: model.TodoResponse[[]model.Todo]{
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
			wantBody: model.TodoResponse[[]model.Todo]{
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
			req := createTestRequest(t, http.MethodPost, "/todos", c.inputBody)

			handler.CreateTodo(rec, req)

			checkMockExpectations(t, mock)
			checkStatusCode(t, c.wantStatusCode, rec.Code)
			got := decodeResponseBody[model.TodoResponse[[]model.Todo]](t, rec)
			checkResponseBody(t, c.wantBody, got)
		})
	}
}

func setUpMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("モックDBの作成に失敗しました: %s", err)
	}
	database.SetDB(db)

	return db, mock
}

func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("満たされていない期待値があります: %s", err)
	}
}

func checkStatusCode(t *testing.T, want, got int) {
	t.Helper()

	if want != got {
		t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", want, got)
	}
}

func decodeResponseBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()

	var got T
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("レスポンスのデコードに失敗しました: %s\nレスポンスボディ: %s", err, rec.Body.String())
	}

	return got
}

func checkResponseBody(t *testing.T, want, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("期待したレスポンス: %v, 実際のレスポンス: %v", want, got)
	}
}

func createTestRequest(t *testing.T, method, path, body string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}

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
