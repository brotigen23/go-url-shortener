package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func MigratePostgresUp(db *sql.DB) error {
	if err := goose.Run("up", db, "./internal/db/migration"); err != nil {
		return err
	}
	return nil
}
func MigratePostgresDown(db *sql.DB) error {
	if err := goose.Run("down", db, "./internal/db/migration"); err != nil {
		return err
	}
	return nil
}
func Reset(db *sql.DB) error {
	if err := goose.Run("reset", db, "./internal/db/migration"); err != nil {
		return err
	}
	return nil
}
