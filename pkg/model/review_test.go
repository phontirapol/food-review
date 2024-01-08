package model_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"food-review/pkg/model"
)

func TestGetAllReviews(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	statement := "SELECT review_id, review FROM review"

	t.Run("No Review", func(t *testing.T) {
		mock.ExpectQuery(statement).
			WillReturnError(sql.ErrNoRows)

		reviews, err := model.GetAllReviews(db)
		if assert.ErrorIs(t, err, sql.ErrNoRows) {
			assert.Nil(t, reviews)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		var fakeReviewID0 uint = 555
		var fakeReviewContent0 string = "This restaurant is overrated"
		var fakeReviewID1 uint = 666
		var fakeReviewContent1 string = "That restaurant is underrated"

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow(fakeReviewID0, fakeReviewContent0).
			AddRow(fakeReviewID1, fakeReviewContent1)

		mock.ExpectQuery(statement).
			WillReturnRows(mockRow)

		reviews, err := model.GetAllReviews(db)
		if assert.NoError(t, err) {
			assert.Len(t, reviews, 2)
			assert.Equal(t, reviews[0].ID, fakeReviewID0)
			assert.Equal(t, reviews[1].Content, fakeReviewContent1)
		}
	})
}

func TestGetReview(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	statement := "SELECT review_id, review FROM review WHERE review_id = ?"

	t.Run("No Review Found", func(t *testing.T) {
		var id uint = 0

		mock.ExpectQuery(statement).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		review, err := model.GetReview(db, id)
		if assert.Error(t, err) {
			assert.Nil(t, review)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		var id uint = 1
		var content string = "Gordan Ramsey is crying"

		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow(id, content)

		mock.ExpectQuery(statement).
			WithArgs(id).
			WillReturnRows(mockRow)

		review, err := model.GetReview(db, id)
		if assert.NoError(t, err) {
			assert.Equal(t, content, review.Content)
		}
	})
}

func TestGetReviewsByKeyword(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	statementFirstHalf := "SELECT review_id, review FROM review WHERE review LIKE '%"
	statementSecondHalf := "%'"

	t.Run("Some DB Error", func(t *testing.T) {
		keyword := "'"
		mock.ExpectQuery(statementFirstHalf + keyword + statementSecondHalf).
			WillReturnError(errors.New("Unknown error"))

		reviews, err := model.GetReviewsByKeyword(db, keyword)
		if assert.Error(t, err) {
			assert.Nil(t, reviews)
			assert.ErrorContains(t, err, "Unknown")
		}
	})

	t.Run("No Review Found", func(t *testing.T) {
		keyword := "cockroach"
		mockRow := sqlmock.NewRows([]string{"review_id", "review"})
		mock.ExpectQuery(statementFirstHalf + keyword + statementSecondHalf).
			WillReturnRows(mockRow)

		reviews, err := model.GetReviewsByKeyword(db, keyword)
		if assert.Error(t, err) {
			assert.Nil(t, reviews)
			assert.ErrorIs(t, err, sql.ErrNoRows)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		keyword := "tiramisu"
		mockRow := sqlmock.NewRows([]string{"review_id", "review"}).
			AddRow(uint(11111), "Worst tiramisu").
			AddRow(uint(22222), "Best tiramisu")
		mock.ExpectQuery(statementFirstHalf + keyword + statementSecondHalf).
			WillReturnRows(mockRow)

		reviews, err := model.GetReviewsByKeyword(db, keyword)
		if assert.NoError(t, err) {
			assert.Len(t, reviews, 2)
			assert.Equal(t, reviews[0].ID, uint(11111))
			assert.Equal(t, reviews[1].Content, "Best tiramisu")
		}
	})
}

func TestKeywordExists(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	statement := "SELECT keyword FROM dictionary WHERE keyword = ?"

	t.Run("Keyword Not Present in Dict", func(t *testing.T) {
		keyword := "tasty cockroach"

		mock.ExpectQuery(statement).
			WillReturnError(sql.ErrNoRows)

		exist, err := model.KeywordExists(db, keyword)

		if assert.NoError(t, err) {
			assert.False(t, exist)
		}
	})

	t.Run("Other Error", func(t *testing.T) {
		keyword := "duck ass"
		errMsg := "Non-sql.ErrNoRows error"

		mock.ExpectQuery(statement).
			WillReturnError(errors.New(errMsg))

		exist, err := model.KeywordExists(db, keyword)

		if assert.Error(t, err) {
			assert.EqualError(t, err, errMsg)
			assert.False(t, exist)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		keyword := "matcha hojicha"

		mockRow := sqlmock.NewRows([]string{"keyword"}).AddRow("mathca hojicha")
		mock.ExpectQuery(statement).
			WillReturnRows(mockRow)

		exist, err := model.KeywordExists(db, keyword)

		if assert.NoError(t, err) {
			assert.True(t, exist)
		}
	})
}

func TestUpdateReview(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	t.Run("JSON Incorrect Data Type", func(t *testing.T) {
		var reviewID uint = 1
		reviewBody := []byte(`{"review": true}`)

		err := model.UpdateReview(db, reviewID, reviewBody)
		expectedError := "cannot unmarshal bool"

		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedError)
		}
	})

	t.Run("Transaction Begin Problem", func(t *testing.T) {
		var reviewID uint = 1
		reviewBody := []byte(`{"review": "This restaurant sucks"}`)

		err := model.UpdateReview(db, reviewID, reviewBody)
		expectedError := "call to database transaction Begin was not expected"

		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedError)
		}
	})

	t.Run("Transaction Prepare Problem", func(t *testing.T) {
		var reviewID uint = 1
		reviewBody := []byte(`{"review": "This restaurant sucks"}`)

		mock.ExpectBegin()

		err := model.UpdateReview(db, reviewID, reviewBody)
		expectedError := "query was not expected"

		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedError)
		}
	})

	t.Run("Transaction Exec Problem", func(t *testing.T) {
		var reviewID uint = 1
		reviewBody := []byte(`{"review": "This restaurant sucks"}`)

		mock.ExpectBegin()

		mock.ExpectPrepare("UPDATE review SET review = ? WHERE review_id = ?")

		err := model.UpdateReview(db, reviewID, reviewBody)
		expectedError := "was not expected"

		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedError)
		}
	})

	t.Run("Transaction Commit Problem", func(t *testing.T) {
		var reviewID uint = 1
		content := "This restaurant sucks"
		reviewBody := []byte(`{"review": "This restaurant sucks"}`)

		mock.ExpectBegin()

		mock.ExpectPrepare("UPDATE review SET review = ? WHERE review_id = ?").
			ExpectExec().
			WithArgs(content, reviewID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := model.UpdateReview(db, reviewID, reviewBody)
		expectedError := "was not expected"

		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedError)
		}
	})

	t.Run("Happy Path", func(t *testing.T) {
		var reviewID uint = 1
		content := "This restaurant sucks"
		reviewBody := []byte(`{"review": "This restaurant sucks"}`)

		mock.ExpectBegin()

		mock.ExpectPrepare("UPDATE review SET review = ? WHERE review_id = ?").
			ExpectExec().
			WithArgs(content, reviewID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit().
			WillReturnError(nil)

		err := model.UpdateReview(db, reviewID, reviewBody)

		if assert.NoError(t, err) {
			assert.Nil(t, err)
		}
	})
}
