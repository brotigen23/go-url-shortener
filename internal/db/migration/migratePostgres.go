package migration

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
)

const (
	TABLE_NAME = "Short_URLs"
)

func MigratePostgres(db *sql.DB) {
	if err := goose.Run("down", db, "./internal/db/migration"); err != nil {
		log.Fatalf("goose error: %v", err)
	}
	if err := goose.Run("up", db, "./internal/db/migration"); err != nil {
		log.Fatalf("goose error: %v", err)
	}
}
