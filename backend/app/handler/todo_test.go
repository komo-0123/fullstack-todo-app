package handler_test

import (
	"backend/app/handler"
	"backend/app/model"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTodos(t *testing.T) {
	cases := map[string]struct {
		mockSetup      func(mock sqlmock.Sqlmock)
		wantStatusCode int
		wantBody       interface{}
	}{
		"正常系": {
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow(1, "title1", false).
						AddRow(2, "title2", true))
			},
			wantStatusCode: http.StatusOK,
			wantBody: createTodoResponse(
				t,
				[]model.Todo{{ID: 1, Title: "title1", IsComplete: false}, {ID: 2, Title: "title2", IsComplete: true}},
				http.StatusOK,
				"",
			),
		},
		"クエリ失敗": {
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnError(fmt.Errorf("DBエラー"))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				[]model.Todo{},
				http.StatusInternalServerError,
				"TODOの取得に失敗しました。",
			),
		},
		"行スキャン失敗": {
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos").
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow("不正なID", "title1", false))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				[]model.Todo{},
				http.StatusInternalServerError,
				"TODOの読み込みに失敗しました。",
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := setUpMockDB(t)
			defer db.Close()

			c.mockSetup(mock)

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
	cases := map[string]struct {
		inputBody      string
		mockSetup      func(mock sqlmock.Sqlmock)
		wantStatusCode int
		wantBody       interface{}
	}{
		"正常系": {
			inputBody: `{"title": "新しいタスク", "is_complete": false}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO todos`).
					WithArgs("新しいタスク", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantStatusCode: http.StatusCreated,
			wantBody: createTodoResponse(
				t,
				[]model.Todo{},
				http.StatusCreated,
				"",
			),
		},
		"不正な入力": {
			inputBody:      `{"title": 123, "is_complete": false}`,
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody: createTodoResponse(
				t,
				[]model.Todo{},
				http.StatusBadRequest,
				"入力が不正です。",
			),
		},
		"DB追加失敗": {
			inputBody: `{"title": "新しいタスク", "is_complete": false}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO todos`).
					WithArgs("新しいタスク", false).
					WillReturnError(fmt.Errorf("DBエラー"))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				[]model.Todo{},
				http.StatusInternalServerError,
				"TODOの追加に失敗しました。",
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := setUpMockDB(t)
			defer db.Close()

			c.mockSetup(mock)

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
