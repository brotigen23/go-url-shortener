package repository

import (
	"database/sql"

	"github.com/brotigen23/go-url-shortener/internal/db/migration"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresRepository(driver string, stringConnection string, logger *zap.Logger) (*PostgresRepository, error) {
	ret := &PostgresRepository{
		logger: logger,
	}
	db, err := sql.Open(driver, stringConnection)
	if err != nil {
		return nil, err
	}
	ret.db = db
	return ret, nil
}

func (r PostgresRepository) CheckDBConnection() error { return r.db.Ping() }
func (r PostgresRepository) Close() error             { return r.db.Close() }

func Migrate(r *PostgresRepository) error {
	return migration.MigratePostgresUp(r.db)
}
func Reset(r *PostgresRepository) error {
	return migration.Reset(r.db)
}
