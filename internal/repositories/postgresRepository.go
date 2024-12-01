package repositories

import (
	"database/sql"
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/db/migration"
	"github.com/brotigen23/go-url-shortener/internal/model"

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

// ------------------------- DEPRECATED ------------------------------------

type P struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresRepository1(driver string, stringConnection string, logger *zap.Logger) (*PostgresRepository, error) {
	ret := &PostgresRepository{
		logger: logger,
	}

	db, err := sql.Open(driver, stringConnection)
	if err != nil {
		return nil, err
	}
	ret.db = db
	migration.MigratePostgresUp(db)
	return ret, nil
}

func (repo *P) CloseConnection() {
	repo.db.Close()
}

func (repo *P) GetByAlias(alias string) (*model.ShortURL, error) {
	query := repo.db.QueryRow(`SELECT * FROM Short_URLs WHERE Alias = $1`, alias)
	var URL string
	var Alias string
	err := query.Scan(&URL, &Alias)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.ShortURL{URL: URL, Alias: Alias}, nil
}
func (repo *P) GetByURL(url string) (*model.ShortURL, error) {
	query := repo.db.QueryRow(`SELECT * FROM Short_URLs WHERE URL = $1`, url)
	var URL string
	var Alias string
	err := query.Scan(&URL, &Alias)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.ShortURL{URL: URL, Alias: Alias}, nil

}

func (repo *P) GetAll() *[]model.ShortURL { return nil }
func (repo *P) Save(model model.ShortURL) error {
	q, err := repo.db.Query("EXPLAIN ANALYZE INSERT INTO Short_URLs(URL, Alias) VALUES($1, $2)", model.URL, model.Alias)
	if err != nil {
		return err
	}
	var plan string
	for q.Next() {
		var s string
		if err := q.Scan(&s); err != nil {
			return err
		}
		plan += s
	}
	repo.logger.Sugar().Infoln(
		"query plan", plan,
	)
	if err := q.Err(); err != nil {
		return err
	}
	return nil
}
func (repo *P) Migrate(model []model.ShortURL) {}
func (repo *P) Close()                         {}
func (repo *P) CheckDBConnection() error {
	return repo.db.Ping()
}

func (repo *P) SaveUser(userID string) error {
	query := "INSERT INTO Users(Name) VALUES($1)"
	_, err := repo.db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *P) GetUserByID(userID string) error {
	query := repo.db.QueryRow(`SELECT * FROM Users WHERE ID = $1`, userID)
	return query.Err()
}

func (repo *P) SaveUserURL(userName string, alias string) error {
	query := "INSERT INTO Users_URLs(User_ID, URL_ID) VALUES((SELECT ID FROM Users WHERE Name = $1), (SELECT ID FROM Short_URLs WHERE Alias = $2))"
	q, err := repo.db.Query(query, userName, alias)
	if err != nil {
		return err
	}
	if err := q.Err(); err != nil {
		return err
	}
	return nil
}

func (repo *P) GetUserURL(userName string) ([]model.ShortURL, error) {
	query := "SELECT url, alias FROM Short_URLs WHERE ID IN (( SELECT URL_ID FROM Users_URLs WHERE User_ID = (SELECT ID FROM Users WHERE Name = $1)))"
	q, err := repo.db.Query(query, userName)
	if err != nil {
		return nil, err
	}
	var id int
	var URL string
	var Alias string
	ret := make([]model.ShortURL, 0)
	for q.Next() {
		err = q.Scan(&id, &URL, &Alias)
		ret = append(ret, *model.NewShortURL(id, URL, Alias))
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}
