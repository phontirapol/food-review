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

func GetReview(db *sql.DB, reviewID uint) (*Review, error) {
	review := Review{}

	statement := "SELECT review_id, review FROM review WHERE review_id = ?"
	row := db.QueryRow(statement, reviewID)
	err := row.Scan(&review.ID, &review.Content)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &review, nil
}
