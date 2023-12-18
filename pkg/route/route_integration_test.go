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

type mockDictionaryDB struct {
	Database *sql.DB
}

func (md *mockDictionaryDB) GetDB() *sql.DB {
	return md.Database
}

func genStatementRev(keyword string) string {
	return "SELECT review_id, review FROM review WHERE review LIKE '%" + keyword + "%'"
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

func TestGetReviewsByKeywordIntegrationService(t *testing.T) {
	dbDict, mockDict, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}
	statementDict := "SELECT keyword FROM dictionary WHERE keyword = ?"

	dbRev, mockRev, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	t.Run("Keyword Not Present", func(t *testing.T) {
		mockDict.ExpectQuery(statementDict).
			WillReturnError(sql.ErrNoRows)

		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDictDB := &mockDictionaryDB{Database: dbDict}
		var mockDictionaryDBOpener db.DictionaryDBOpener = mockDictDB

		mockHandler := &route.Handler{
			Template:     mockTemplater,
			DictionaryDB: mockDictionaryDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews?query=foodkeyword", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReviewsByKeyword)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusUnprocessableEntity
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Other Error - Keyword Exist Section", func(t *testing.T) {
		mockDict.ExpectQuery(statementDict).
			WillReturnError(errors.New("Some other error in keyword db"))

		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDictDB := &mockDictionaryDB{Database: dbDict}
		var mockDictionaryDBOpener db.DictionaryDBOpener = mockDictDB

		mockHandler := &route.Handler{
			Template:     mockTemplater,
			DictionaryDB: mockDictionaryDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews?query=foodkeyword", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReviewsByKeyword)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusInternalServerError
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("No Review with this Keyword", func(t *testing.T) {
		foodKeyword := "foie gras"
		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)

		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDictDB := &mockDictionaryDB{Database: dbDict}
		var mockDictionaryDBOpener db.DictionaryDBOpener = mockDictDB

		mockRevStatement := genStatementRev(foodKeyword)
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnError(sql.ErrNoRows)

		mockRevDB := &mockReviewDB{Database: dbRev}
		var mockReviewDBOpener db.ReviewDBOpener = mockRevDB

		mockHandler := &route.Handler{
			Template:     mockTemplater,
			ReviewDB:     mockReviewDBOpener,
			DictionaryDB: mockDictionaryDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews?query="+foodKeyword, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReviewsByKeyword)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusUnprocessableEntity
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Other Error - Get Review Section", func(t *testing.T) {
		foodKeyword := "foie gras"
		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)

		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDictDB := &mockDictionaryDB{Database: dbDict}
		var mockDictionaryDBOpener db.DictionaryDBOpener = mockDictDB

		mockRevStatement := genStatementRev(foodKeyword)
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnError(errors.New("Some other error in review db"))

		mockRevDB := &mockReviewDB{Database: dbRev}
		var mockReviewDBOpener db.ReviewDBOpener = mockRevDB

		mockHandler := &route.Handler{
			Template:     mockTemplater,
			ReviewDB:     mockReviewDBOpener,
			DictionaryDB: mockDictionaryDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews?query="+foodKeyword, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReviewsByKeyword)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusInternalServerError
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Template Error", func(t *testing.T) {
		foodKeyword := "foie gras"
		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)

		mockTmpl := &mockTemplate{errMsg: errors.New("Some error in template")}
		var mockTemplater template.Templater = mockTmpl

		mockDictDB := &mockDictionaryDB{Database: dbDict}
		var mockDictionaryDBOpener db.DictionaryDBOpener = mockDictDB

		mockRevStatement := genStatementRev(foodKeyword)
		mockRevRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("8888", "this restaurant sucks")
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnRows(mockRevRow)

		mockRevDB := &mockReviewDB{Database: dbRev}
		var mockReviewDBOpener db.ReviewDBOpener = mockRevDB

		mockHandler := &route.Handler{
			Template:     mockTemplater,
			ReviewDB:     mockReviewDBOpener,
			DictionaryDB: mockDictionaryDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews?query="+foodKeyword, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReviewsByKeyword)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusInternalServerError
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		foodKeyword := "foie gras"
		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)

		mockTmpl := &mockTemplate{errMsg: nil}
		var mockTemplater template.Templater = mockTmpl

		mockDictDB := &mockDictionaryDB{Database: dbDict}
		var mockDictionaryDBOpener db.DictionaryDBOpener = mockDictDB

		mockRevStatement := genStatementRev(foodKeyword)
		mockRevRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("8888", "this restaurant sucks")
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnRows(mockRevRow)

		mockRevDB := &mockReviewDB{Database: dbRev}
		var mockReviewDBOpener db.ReviewDBOpener = mockRevDB

		mockHandler := &route.Handler{
			Template:     mockTemplater,
			ReviewDB:     mockReviewDBOpener,
			DictionaryDB: mockDictionaryDBOpener,
		}

		r, err := http.NewRequest(http.MethodGet, "/reviews?query="+foodKeyword, nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetReviewsByKeyword)
		handler.ServeHTTP(w, r)

		expectedStatus := http.StatusOK
		if status := w.Code; status != expectedStatus {
			t.Errorf("expected %v code but got %v instead", expectedStatus, status)
		}
	})
}
