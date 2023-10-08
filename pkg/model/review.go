package model

import (
	"database/sql"

	"food-review/pkg/db"
)

type Review struct {
	ID      uint   `json:"review_id"`
	Content string `json:"review"`
}

func GetAllReviews() ([]Review, error) {
	var allReviews []Review

	db := db.InitReviewDB()

	statement := "SELECT review_id, review FROM review"
	rows, err := db.Database.Query(statement)
	if err == sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		review := Review{}
		err = rows.Scan(
			&review.ID,
			&review.Content,
		)
		if err != nil {
			return nil, err
		}

		allReviews = append(allReviews, review)
	}

	return allReviews, nil
}
