package route_test

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"food-review/pkg/db"
	"food-review/pkg/route"
	"food-review/pkg/template"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

type mockTemplate struct {
	errMsg error
}

func (mt *mockTemplate) ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error {
	return mt.errMsg
}

type mockReviewDB struct {
	Database *sql.DB
}

func (mr *mockReviewDB) GetDB() *sql.DB {
	return mr.Database
}

func TestGetAllReviewsIntegrationService(t *testing.T) {
	t.Run("No Review in Database", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error not expected while opening mock db, %v", err)
		}

		statement := "SELECT review_id, review FROM review"
		mock.ExpectQuery(statement).WillReturnError(sql.ErrNoRows)

		mockRDB := &mockReviewDB{Database: mockDB}
		var mockReviewDBOpener db.ReviewDBOpener = mockRDB

		mockHandler := &route.Handler{
			Template: mockTemplater,
			ReviewDB: mockReviewDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetAllReviews)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusInternalServerError
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Error Caused in Template", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: errors.New("Some template error")}
		var mockTemplater template.Templater = mockTmpl

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error not expected while opening mock db, %v", err)
		}

		statement := "SELECT review_id, review FROM review"
		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("999999", "Integration review")
		mock.ExpectQuery(statement).
			WillReturnRows(mockRow)

		mockRDB := &mockReviewDB{Database: mockDB}
		var mockReviewDBOpener db.ReviewDBOpener = mockRDB

		mockHandler := &route.Handler{
			Template: mockTemplater,
			ReviewDB: mockReviewDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetAllReviews)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusInternalServerError
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error not expected while opening mock db, %v", err)
		}

		statement := "SELECT review_id, review FROM review"
		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("999999", "Integration review")
		mock.ExpectQuery(statement).
			WillReturnRows(mockRow)

		mockRDB := &mockReviewDB{Database: mockDB}
		var mockReviewDBOpener db.ReviewDBOpener = mockRDB

		mockHandler := &route.Handler{
			Template: mockTemplater,
			ReviewDB: mockReviewDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetAllReviews)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusOK
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})
}

func TestGetReviewIntegrationService(t *testing.T) {
	t.Run("Invalid ID", func(t *testing.T) {
		testSuite := []string{
			"-1",
			"1.2",
			"abc",
			"1/1",
		}

		for _, testCase := range testSuite {
			mockTmpl := &mockTemplate{errMsg: nil}
			var mockTemplater template.Templater = mockTmpl

			mockRDB := &mockReviewDB{}
			var mockReviewDBOpener db.ReviewDBOpener = mockRDB

			mockHandler := &route.Handler{
				Template: mockTemplater,
				ReviewDB: mockReviewDBOpener,
			}

			r, err := http.NewRequest(http.MethodGet, "/reviews/"+testCase, nil)
			if err != nil {
				t.Error(err)
			}

			vars := map[string]string{"reviewID": testCase}
			r = mux.SetURLVars(r, vars)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(mockHandler.GetReview)
			handler.ServeHTTP(w, r)

			expectedStatus := http.StatusBadRequest
			if status := w.Code; status != expectedStatus {
				t.Errorf("expected %v code but got %v instead", expectedStatus, status)
			}
		}
	})

	t.Run("No Review with this ID", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error not expected while opening mock db, %v", err)
		}

		statement := "SELECT review_id, review FROM review WHERE review_id = ?"
		mock.ExpectQuery(statement).WillReturnError(sql.ErrNoRows)

		mockRDB := &mockReviewDB{Database: mockDB}
		var mockReviewDBOpener db.ReviewDBOpener = mockRDB

		mockHandler := &route.Handler{
			Template: mockTemplater,
			ReviewDB: mockReviewDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews/9999999", nil)
		if err != nil {
			t.Error(err)
		}

		vars := map[string]string{"reviewID": "9999999"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReview)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusUnprocessableEntity
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Error in Template", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: errors.New("Some error in template")}
		var mockTemplater template.Templater = mockTmpl

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error not expected while opening mock db, %v", err)
		}

		statement := "SELECT review_id, review FROM review WHERE review_id = ?"
		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("1", "This restaurant deserves 9 Michelin stars")
		mock.ExpectQuery(statement).WillReturnRows(mockRow)

		mockRDB := &mockReviewDB{Database: mockDB}
		var mockReviewDBOpener db.ReviewDBOpener = mockRDB

		mockHandler := &route.Handler{
			Template: mockTemplater,
			ReviewDB: mockReviewDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews/1", nil)
		if err != nil {
			t.Error(err)
		}

		vars := map[string]string{"reviewID": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReview)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusInternalServerError
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("error not expected while opening mock db, %v", err)
		}

		statement := "SELECT review_id, review FROM review WHERE review_id = ?"
		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("1", "This restaurant deserves 9 Michelin stars")
		mock.ExpectQuery(statement).WillReturnRows(mockRow)

		mockRDB := &mockReviewDB{Database: mockDB}
		var mockReviewDBOpener db.ReviewDBOpener = mockRDB

		mockHandler := &route.Handler{
			Template: mockTemplater,
			ReviewDB: mockReviewDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews/1", nil)
		if err != nil {
			t.Error(err)
		}

		vars := map[string]string{"reviewID": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReview)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusOK
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})
}
