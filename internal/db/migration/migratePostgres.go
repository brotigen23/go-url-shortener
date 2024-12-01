package migration

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
)

func MigratePostgresUp(db *sql.DB) {
	if err := goose.Run("up", db, "./internal/db/migration"); err != nil {
		log.Fatalf("goose error: %v", err)
	}
}
func MigratePostgresDown(db *sql.DB) {
	if err := goose.Run("down", db, "./internal/db/migration"); err != nil {
		log.Fatalf("goose error: %v", err)
	}
}
