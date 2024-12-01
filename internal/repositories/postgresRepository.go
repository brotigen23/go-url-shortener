package repositories

import (
	"database/sql"

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
	//migration.MigratePostgres(db)
	return ret, nil
}

func (repo PostgresRepository) CheckDBConnection() error { return repo.db.Ping() }
func (repo PostgresRepository) Close() error             { return repo.db.Close() }
