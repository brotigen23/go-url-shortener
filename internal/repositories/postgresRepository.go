package repositories

import (
	"database/sql"

	"github.com/brotigen23/go-url-shortener/internal/model"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(stringConnection string) *PostgresRepository {
	ret := &PostgresRepository{}
	db, err := sql.Open("postgres", stringConnection)
	if err != nil {
		panic(err)
	}
	ret.db = db

	return ret
}

func (repo *PostgresRepository) CloseConnection() {
	repo.db.Close()
}

func (repo *PostgresRepository) GetByAlias(alias string) (*model.Alias, error) { return nil, nil }
func (repo *PostgresRepository) GetByURL(url string) (*model.Alias, error)     { return nil, nil }
func (repo *PostgresRepository) GetAll() *[]model.Alias                        { return nil }
func (repo *PostgresRepository) Save(model model.Alias) error                  { return nil }
func (repo *PostgresRepository) Migrate(model []model.Alias)                   {}
func (repo *PostgresRepository) Close()                                        {}
func (repo *PostgresRepository) CheckDBConnection() error {
	return repo.db.Ping()
}
