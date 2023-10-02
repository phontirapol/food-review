package db_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"food-review/pkg/db"
)

func TestInit(t *testing.T) {
	t.Run("Invalid SQL Driver", func(t *testing.T) {
		review := &db.ReviewDB{
			Driver: "unknownsql",
		}

		expectedErrorSubstring := "sql: unknown driver"
		err := review.Init()
		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedErrorSubstring)
		}
	})

	t.Run("Invalid SQL Init Statement", func(t *testing.T) {
		review := &db.ReviewDB{
			Driver: "sqlite3",
			InitStatement: `
				CREATE TABLE IF NOT EXIT
				review (
					review_id INTEGER PRIMARY KEY,
					review TEXT
				)
			`,
		}

		expectedErrorSubstring := "syntax error"
		err := review.Init()
		if assert.Error(t, err) {
			assert.ErrorContains(t, err, expectedErrorSubstring)
		}
	})

	t.Run("Successful DB Initialization", func(t *testing.T) {
		review := &db.ReviewDB{
			Driver:     "sqlite3",
			DataSource: "./reviewtest.db",
			InitStatement: `
				CREATE TABLE IF NOT EXISTS
				review (
					review_id INTEGER PRIMARY KEY,
					review TEXT
				)
			`,
		}

		err := review.Init()
		if assert.NoError(t, err) {
			if err := os.Remove("reviewtest.db"); err != nil {
				t.Error(err)
			}
		}
	})
}
