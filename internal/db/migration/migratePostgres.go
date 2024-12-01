package migration

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
)

func MigratePostgres(db *sql.DB) {
	if err := goose.Run("up", db, "./internal/db/migration", "up"); err != nil {
		log.Fatalf("goose error: %v", err)
	}
}
