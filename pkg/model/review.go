package model

import (
	"database/sql"
	"encoding/json"
)

type Review struct {
	ID      uint   `json:"review_id"`
	Content string `json:"review"`
	Keyword string
}

func GetAllReviews(db *sql.DB) ([]*Review, error) {
	var allReviews []*Review

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
		allReviews = append(allReviews, &review)
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

func GetReviewsByKeyword(db *sql.DB, keyword string) ([]*Review, error) {
	var targetReviews []*Review

	statement := "SELECT review_id, review FROM review WHERE review LIKE '%" + keyword + "%'"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		review := Review{}
		_ = rows.Scan(
			&review.ID,
			&review.Content,
		)
		review.Keyword = keyword

		targetReviews = append(targetReviews, &review)
	}

	if len(targetReviews) == 0 {
		return nil, sql.ErrNoRows
	}

	return targetReviews, nil
}

func KeywordExists(db *sql.DB, keyword string) (bool, error) {
	var foodKey string
	statement := "SELECT keyword FROM dictionary WHERE keyword = ?"
	row := db.QueryRow(statement, keyword)
	err := row.Scan(&foodKey)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateReview(db *sql.DB, reviewID uint, reviewBody []byte) error {
	editedReview := Review{ID: reviewID}

	err := json.Unmarshal(reviewBody, &editedReview)
	if err != nil {
		return err
	}

	ps, err := db.Begin()
	if err != nil {
		return err
	}

	updateStatement := "UPDATE review SET review = ? WHERE review_id = ?"
	statement, err := db.Prepare(updateStatement)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(editedReview.Content, editedReview.ID)
	if err != nil {
		return err
	}

	err = ps.Commit()
	if err != nil {
		return err
	}

	return nil
}
