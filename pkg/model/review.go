package model

import (
	"database/sql"
)

type Review struct {
	ID      uint   `json:"review_id"`
	Content string `json:"review"`
}

func GetAllReviews(db *sql.DB) ([]Review, error) {
	var allReviews []Review

	statement := "SELECT review_id, review FROM review"
	rows, err := db.Query(statement)
	if err == sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		review := Review{}
		_ = rows.Scan(
			&review.ID,
			&review.Content,
		)
		allReviews = append(allReviews, review)
	}

	return allReviews, nil
}
