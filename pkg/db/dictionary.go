package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDictionaryDB() *DictionaryDB {
	driver := "sqlite3"
	dataSource := "./db/dictionary.db"
	initStatement := `
		CREATE TABLE IF NOT EXISTS
		dictionary (
			keyword TEXT
		)
	`
	db := &DictionaryDB{
		Driver:        driver,
		DataSource:    dataSource,
		InitStatement: initStatement,
	}

	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type DictionaryDBOpener interface {
	GetDB() *sql.DB
}

type DictionaryDB struct {
	Driver        string
	DataSource    string
	InitStatement string
	Database      *sql.DB
}

func (db *DictionaryDB) Init() error {
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

func (db *DictionaryDB) GetDB() *sql.DB {
	return db.Database
}
