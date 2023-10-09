package model_test

import (
	"database/sql"
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
