package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

func MigratePostgresUp(db *sql.DB) error {
	if err := goose.Run("up", db, "./internal/database/migration"); err != nil {
		return err
	}
	return nil
}
func MigratePostgresDown(db *sql.DB) error {
	if err := goose.Run("down", db, "../database/migration"); err != nil {
		return err
	}
	return nil
}
func Reset(db *sql.DB) error {
	if err := goose.Run("reset", db, "../database/migration"); err != nil {
		return err
	}
	return nil
}
