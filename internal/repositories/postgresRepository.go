package repositories

import (
	"database/sql"
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/db/migration"
	"github.com/brotigen23/go-url-shortener/internal/model"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(stringConnection string) (*PostgresRepository, error) {
	ret := &PostgresRepository{}

	db, err := sql.Open("postgres", stringConnection)
	if err != nil {
		return nil, err
	}
	ret.db = db
	migration.MigratePostgres(db)
	return ret, nil
}

func (repo *PostgresRepository) CloseConnection() {
	repo.db.Close()
}

func (repo *PostgresRepository) GetByAlias(alias string) (*model.Alias, error) {
	query := repo.db.QueryRow(`SELECT * FROM Aliases WHERE "Alias" = $1`, alias)
	var URL string
	var Alias string
	err := query.Scan(&URL, &Alias)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.Alias{URL: URL, Alias: Alias}, nil
}
func (repo *PostgresRepository) GetByURL(url string) (*model.Alias, error) {
	query := repo.db.QueryRow(`SELECT * FROM Aliases WHERE "URL" = $1`, url)
	var URL string
	var Alias string
	err := query.Scan(&URL, &Alias)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.Alias{URL: URL, Alias: Alias}, nil

}

func (repo *PostgresRepository) GetAll() *[]model.Alias { return nil }
func (repo *PostgresRepository) Save(model model.Alias) error {
	result, err := repo.db.Exec("INSERT INTO Aliases VALUES($1, $2)", model.URL, model.Alias)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
func (repo *PostgresRepository) Migrate(model []model.Alias) {}
func (repo *PostgresRepository) Close()                      {}
func (repo *PostgresRepository) CheckDBConnection() error {
	return repo.db.Ping()
}
