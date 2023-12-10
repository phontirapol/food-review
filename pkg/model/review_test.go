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
