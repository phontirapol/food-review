package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitReviewDB() {
	driver := "sqlite3"
	dataSource := "./review.db"
	initStatement := `
		CREATE TABLE IF NOT EXISTS
		review (
			review_id INTEGER PRIMARY KEY,
			review TEXT
		)
	`
	db := &ReviewDB{
		Driver:        driver,
		DataSource:    dataSource,
		InitStatement: initStatement,
	}

	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
}

type ReviewDB struct {
	Driver        string
	DataSource    string
	InitStatement string
	Database      *sql.DB
}

func (db *ReviewDB) Init() error {
	database, err := sql.Open(db.Driver, db.DataSource)
	if err != nil {
		return err
	}

	db.Database = database
	_, err = db.Database.Exec(db.InitStatement)
	if err != nil {
		return err
	}

	return nil
}
