package migration

import (
	"database/sql"

	"github.com/pressly/goose"
)

// Производит миграцию базы данных
func MigratePostgresUp(db *sql.DB) error {
	if err := goose.Run("up", db, "./internal/database/migration"); err != nil {
		return err
	}
	return nil
}
