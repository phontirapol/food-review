package route_test

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"sync"

	"food-review/pkg/route"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	GET string = http.MethodGet
	PUT string = http.MethodPut
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

func constructHandler(
	templateMock *mockTemplate,
	reviewDBMock *mockReviewDB,
	dictDBMock *mockDictionaryDB,
) *route.Handler {
	return &route.Handler{
		Template:     templateMock,
		ReviewDB:     reviewDBMock,
		DictionaryDB: dictDBMock,
	}
}

func testHandler(
	t *testing.T,
	targetHandler http.HandlerFunc,
	method string,
	url string,
	body io.Reader,
	varObj map[string]string,
	expectedStatus int,
) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Error(err)
	}

	if varObj != nil {
		r = mux.SetURLVars(r, varObj)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(targetHandler)
	handler.ServeHTTP(w, r)

	if w.Code != expectedStatus {
		t.Errorf("expected %v code but got %v instead", expectedStatus, w.Code)
	}
}

func reviewStatement(keyword string) string {
	return "SELECT review_id, review FROM review WHERE review LIKE '%" + keyword + "%'"
}

func TestGetAllReviewsIntegrationService(t *testing.T) {
	url := "/reviews"

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error not expected while opening mock db, %v", err)
	}
	statement := "SELECT review_id, review FROM review"

	t.Run("No Review in Database", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		statement := "SELECT review_id, review FROM review"
		mock.ExpectQuery(statement).WillReturnError(sql.ErrNoRows)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		testHandler(t, mockHandler.GetAllReviews, GET, url, nil, nil, http.StatusInternalServerError)
	})

	t.Run("Error Caused in Template", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: errors.New("Some template error")}

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("999999", "Integration review")
		mock.ExpectQuery(statement).
			WillReturnRows(mockRow)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		testHandler(t, mockHandler.GetAllReviews, GET, url, nil, nil, http.StatusInternalServerError)
	})

	t.Run("Happy Path", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("999999", "Integration review")
		mock.ExpectQuery(statement).
			WillReturnRows(mockRow)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		testHandler(t, mockHandler.GetAllReviews, GET, url, nil, nil, http.StatusOK)
	})
}

func TestGetReviewIntegrationService(t *testing.T) {
	url := "/reviews/"

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error not expected while opening mock db, %v", err)
	}
	statement := "SELECT review_id, review FROM review WHERE review_id = ?"

	t.Run("Invalid ID", func(t *testing.T) {
		testSuite := []string{
			"-1",
			"1.2",
			"abc",
			"1/1",
		}

		for _, testCase := range testSuite {
			mockTmpl := &mockTemplate{errMsg: nil}

			mockRDB := &mockReviewDB{}

			mockHandler := constructHandler(mockTmpl, mockRDB, nil)

			vars := map[string]string{"reviewID": testCase}
			testHandler(t, mockHandler.GetReview, GET, url+testCase, nil, vars, http.StatusBadRequest)
		}
	})

	t.Run("No Review with this ID", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		mock.ExpectQuery(statement).WillReturnError(sql.ErrNoRows)

		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		id := "9999999"
		vars := map[string]string{"reviewID": id}
		testHandler(t, mockHandler.GetReview, GET, url+id, nil, vars, http.StatusUnprocessableEntity)
	})

	t.Run("Error in Template", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: errors.New("Some error in template")}

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("1", "This restaurant deserves 9 Michelin stars")
		mock.ExpectQuery(statement).WillReturnRows(mockRow)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		id := "1"
		vars := map[string]string{"reviewID": id}
		testHandler(t, mockHandler.GetReview, GET, url+id, nil, vars, http.StatusInternalServerError)
	})

	t.Run("Happy Path", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("1", "This restaurant deserves 9 Michelin stars")
		mock.ExpectQuery(statement).WillReturnRows(mockRow)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		id := "1"
		vars := map[string]string{"reviewID": id}
		testHandler(t, mockHandler.GetReview, GET, url+id, nil, vars, http.StatusOK)
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
		mockTmpl := &mockTemplate{errMsg: nil}

		mockDict.ExpectQuery(statementDict).
			WillReturnError(sql.ErrNoRows)
		mockDictDB := &mockDictionaryDB{Database: dbDict}

		mockHandler := constructHandler(mockTmpl, nil, mockDictDB)

		url := "/reviews?query=foodkeyword"
		testHandler(t, mockHandler.GetReviewsByKeyword, GET, url, nil, nil, http.StatusUnprocessableEntity)
	})

	t.Run("Other Error - Keyword Exist Section", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		mockDict.ExpectQuery(statementDict).
			WillReturnError(errors.New("Some other error in keyword db"))
		mockDictDB := &mockDictionaryDB{Database: dbDict}

		mockHandler := constructHandler(mockTmpl, nil, mockDictDB)

		url := "/reviews?query=foodkeyword"
		testHandler(t, mockHandler.GetReviewsByKeyword, GET, url, nil, nil, http.StatusInternalServerError)
	})

	t.Run("No Review with this Keyword", func(t *testing.T) {
		foodKeyword := "foie gras"

		mockTmpl := &mockTemplate{errMsg: nil}

		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)
		mockDictDB := &mockDictionaryDB{Database: dbDict}

		mockRevStatement := reviewStatement(foodKeyword)
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnError(sql.ErrNoRows)
		mockRevDB := &mockReviewDB{Database: dbRev}

		mockHandler := constructHandler(mockTmpl, mockRevDB, mockDictDB)

		url := "/reviews?query=" + foodKeyword
		testHandler(t, mockHandler.GetReviewsByKeyword, GET, url, nil, nil, http.StatusUnprocessableEntity)
	})

	t.Run("Other Error - Get Review Section", func(t *testing.T) {
		foodKeyword := "foie gras"

		mockTmpl := &mockTemplate{errMsg: nil}

		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)
		mockDictDB := &mockDictionaryDB{Database: dbDict}

		mockRevStatement := reviewStatement(foodKeyword)
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnError(errors.New("Some other error in review db"))
		mockRevDB := &mockReviewDB{Database: dbRev}

		mockHandler := constructHandler(mockTmpl, mockRevDB, mockDictDB)

		url := "/reviews?query=" + foodKeyword
		testHandler(t, mockHandler.GetReviewsByKeyword, GET, url, nil, nil, http.StatusInternalServerError)
	})

	t.Run("Template Error", func(t *testing.T) {
		foodKeyword := "foie gras"

		mockTmpl := &mockTemplate{errMsg: errors.New("Some error in template")}

		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)
		mockDictDB := &mockDictionaryDB{Database: dbDict}

		mockRevStatement := reviewStatement(foodKeyword)
		mockRevRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("8888", "this restaurant sucks")
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnRows(mockRevRow)
		mockRevDB := &mockReviewDB{Database: dbRev}

		mockHandler := constructHandler(mockTmpl, mockRevDB, mockDictDB)

		url := "/reviews?query=" + foodKeyword
		testHandler(t, mockHandler.GetReviewsByKeyword, GET, url, nil, nil, http.StatusInternalServerError)
	})

	t.Run("Happy Path", func(t *testing.T) {
		foodKeyword := "foie gras"

		mockTmpl := &mockTemplate{errMsg: nil}

		mockDictRow := sqlmock.NewRows([]string{"keyword"}).AddRow(foodKeyword)
		mockDict.ExpectQuery(statementDict).
			WillReturnRows(mockDictRow)
		mockDictDB := &mockDictionaryDB{Database: dbDict}

		mockRevStatement := reviewStatement(foodKeyword)
		mockRevRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("8888", "this restaurant sucks")
		mockRev.ExpectQuery(mockRevStatement).
			WillReturnRows(mockRevRow)
		mockRevDB := &mockReviewDB{Database: dbRev}

		mockHandler := constructHandler(mockTmpl, mockRevDB, mockDictDB)

		url := "/reviews?query=" + foodKeyword
		testHandler(t, mockHandler.GetReviewsByKeyword, GET, url, nil, nil, http.StatusOK)
	})
}

func TestAccessReviewEditIntegrationService(t *testing.T) {
	url := "/reviews/"
	suffix := "/edit"

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error not expected while opening mock db, %v", err)
	}
	statement := "SELECT review_id, review FROM review WHERE review_id = ?"

	t.Run("Invalid ID", func(t *testing.T) {
		testSuite := []string{
			"-1",
			"1.2",
			"abc",
			"1/1",
		}

		for _, testCase := range testSuite {
			mockTmpl := &mockTemplate{errMsg: nil}

			mockHandler := constructHandler(mockTmpl, nil, nil)

			vars := map[string]string{"reviewID": testCase}
			testHandler(t, mockHandler.AccessReviewEdit, GET, url+testCase+suffix, nil, vars, http.StatusBadRequest)
		}
	})

	t.Run("No Review with this ID", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		mock.ExpectQuery(statement).WillReturnError(sql.ErrNoRows)

		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		id := "9999999"
		vars := map[string]string{"reviewID": id}
		testHandler(t, mockHandler.AccessReviewEdit, GET, url+id+suffix, nil, vars, http.StatusUnprocessableEntity)
	})

	t.Run("Error in Template", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: errors.New("Some error in template")}

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("1", "This restaurant deserves 9 Michelin stars")
		mock.ExpectQuery(statement).WillReturnRows(mockRow)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		id := "1"
		vars := map[string]string{"reviewID": id}
		testHandler(t, mockHandler.AccessReviewEdit, GET, url+id+suffix, nil, vars, http.StatusInternalServerError)
	})

	t.Run("Happy Path", func(t *testing.T) {
		mockTmpl := &mockTemplate{errMsg: nil}

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow("1", "This restaurant deserves 9 Michelin stars")
		mock.ExpectQuery(statement).WillReturnRows(mockRow)
		mockRDB := &mockReviewDB{Database: mockDB}

		mockHandler := constructHandler(mockTmpl, mockRDB, nil)

		id := "1"
		vars := map[string]string{"reviewID": id}
		testHandler(t, mockHandler.AccessReviewEdit, GET, url+id+suffix, nil, vars, http.StatusOK)
	})
}

func TestEditReviewConcurrency(t *testing.T) {
	t.Run("Update Remains Concurrent", func(t *testing.T) {
		dbRev, mockRev, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Error(err)
		}

		mockTmpl := &mockTemplate{errMsg: nil}
		mockDictDB := &mockDictionaryDB{Database: nil}

		mockRev.ExpectBegin()

		mockRev.ExpectPrepare("UPDATE review SET review = ? WHERE review_id = ?").
			ExpectExec().
			WithArgs("This is great", uint(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockRev.ExpectCommit().
			WillReturnError(nil)

		mockRevDB := &mockReviewDB{Database: dbRev}
		mockHandler := constructHandler(mockTmpl, mockRevDB, mockDictDB)

		bodyContents := []string{
			`
			{
				"review_id": 1,
				"review": "This is great"
			}
			`,
			`
			{
				"review_id": 1,
				"review": "This is bad"
			}
			`,
		}

		wg := sync.WaitGroup{}
		mu := sync.Mutex{}
		completionTime := []time.Time{}

		for _, body := range bodyContents {
			wg.Add(1)

			go func(body string) {
				defer wg.Done()

				r, err := http.NewRequest(http.MethodPut, "/reviews", strings.NewReader(body))
				if err != nil {
					t.Error(err)
				}

				vars := map[string]string{"reviewID": "1"}
				r = mux.SetURLVars(r, vars)

				w := httptest.NewRecorder()
				handler := http.HandlerFunc(mockHandler.EditReview)
				handler.ServeHTTP(w, r)

				mu.Lock()
				defer mu.Unlock()
				completionTime = append(completionTime, time.Now())
			}(body)
		}

		wg.Wait()

		assert.True(t, completionTime[1].After(completionTime[0]))
	})
}
