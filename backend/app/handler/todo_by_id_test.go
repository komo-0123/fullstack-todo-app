package handler_test

import (
	"backend/app/handler"
	"backend/app/model"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTodoById(t *testing.T) {
	cases := map[string]struct {
		ID             int
		mockSetup      func(mock sqlmock.Sqlmock)
		wantStatusCode int
		wantBody       interface{}
	}{
		"正常系": {
			ID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow(1, "title1", false))
			},
			wantStatusCode: http.StatusOK,
			wantBody: createTodoResponse(
				t,
				&model.Todo{ID: 1, Title: "title1", IsComplete: false},
				http.StatusOK,
				"",
			),
		},
		"IDが不正": {
			ID:             0,
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusInternalServerError,
				"TODOの取得に失敗しました。",
			),
		},
		"TODOが存在しない": {
			ID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantStatusCode: http.StatusNotFound,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusNotFound,
				"TODOが見つかりません。",
			),
		},
		"クエリ失敗": {
			ID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusInternalServerError,
				"TODOの取得に失敗しました。",
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := setUpMockDB(t)
			defer db.Close()

			c.mockSetup(mock)

			rec := httptest.NewRecorder()
			path := "/todos/" + strconv.Itoa(c.ID)
			req := createTestRequest(t, http.MethodGet, path, "")

			handler.GetTodoById(rec, req)

			checkMockExpectations(t, mock)
			checkStatusCode(t, c.wantStatusCode, rec.Code)
			got := decodeResponseBody[model.TodoResponse](t, rec)
			checkResponseBody(t, c.wantBody, got)
		})
	}
}

func TestUpdateTodoById(t *testing.T) {
	cases := map[string]struct {
		ID             int
		inputBody      string
		mockSetup      func(mock sqlmock.Sqlmock)
		wantStatusCode int
		wantBody       interface{}
	}{
		"正常系": {
			ID:        1,
			inputBody: `{"title": "Updated Title", "is_complete": true}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT id, title, is_complete FROM todos WHERE id = \?$`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_complete"}).
						AddRow(1, "Existing Title", false))
				mock.ExpectExec(`^UPDATE todos SET title = \?, is_complete = \? WHERE id = \?$`).
					WithArgs("Updated Title", true, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantStatusCode: http.StatusOK,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusOK,
				"",
			),
		},
		"TODOが見つかりません": {
			ID:        1,
			inputBody: `{"title": "Updated Title", "is_complete": true}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			wantStatusCode: http.StatusNotFound,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusNotFound,
				"TODOが見つかりません。",
			),
		},
		"クエリ失敗": {
			ID:        1,
			inputBody: `{"title": "Updated Title", "is_complete": true}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, is_complete FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				nil,
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
			path := "/todos/" + strconv.Itoa(c.ID)
			req := createTestRequest(t, http.MethodPut, path, c.inputBody)

			handler.UpdateTodoById(rec, req)

			checkMockExpectations(t, mock)
			checkStatusCode(t, c.wantStatusCode, rec.Code)
			got := decodeResponseBody[model.TodoResponse](t, rec)
			checkResponseBody(t, c.wantBody, got)
		})
	}
}

func TestDeleteTodoById(t *testing.T) {
	cases := map[string]struct {
		ID             int
		mockSetup      func(mock sqlmock.Sqlmock)
		wantStatusCode int
		wantBody       interface{}
	}{
		"正常系": {
			ID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantStatusCode: http.StatusOK,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusOK,
				"",
			),
		},
		"TODOが見つかりません": {
			ID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantStatusCode: http.StatusNotFound,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusNotFound,
				"指定のTODOは削除済みです。",
			),
		},
		"クエリ失敗": {
			ID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM todos WHERE id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: createTodoResponse(
				t,
				nil,
				http.StatusInternalServerError,
				"TODOの削除に失敗しました。",
			),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := setUpMockDB(t)
			defer db.Close()

			c.mockSetup(mock)

			rec := httptest.NewRecorder()
			path := "/todos/" + strconv.Itoa(c.ID)
			req := createTestRequest(t, http.MethodDelete, path, "")

			handler.DeleteTodoById(rec, req)

			checkMockExpectations(t, mock)
			checkStatusCode(t, c.wantStatusCode, rec.Code)
			got := decodeResponseBody[model.TodoResponse](t, rec)
			checkResponseBody(t, c.wantBody, got)
		})
	}
}
