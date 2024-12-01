package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

//---------------------- ShortURLs table ----------------------

func (repo PostgresRepository) GetAllShortURL() ([]model.ShortURL, error) { return nil, nil }

func (repo PostgresRepository) GetShortURLByID(id int) (*model.ShortURL, error) {
	query := repo.db.QueryRow(`SELECT * FROM Short_URLs WHERE ID = $1`, id)
	var ID int
	var URL string
	var Alias string
	err := query.Scan(&ID, &URL, &Alias)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.ShortURL{Id: ID, URL: URL, Alias: Alias}, nil
}
func (repo PostgresRepository) GetShortURLByAlias(alias string) (*model.ShortURL, error) {
	query := repo.db.QueryRow(`SELECT * FROM Short_URLs WHERE Alias = $1`, alias)
	var ID int
	var URL string
	var Alias string
	err := query.Scan(&ID, &URL, &Alias)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.ShortURL{Id: ID, URL: URL, Alias: Alias}, nil
}

func (repo PostgresRepository) GetShortURLByURL(URL string) (*model.ShortURL, error) { return nil, nil }

func (repo PostgresRepository) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	query := "INSERT INTO Short_URLs(URL, Alias) VALUES($1, $2) RETURNING ID"
	var (
		id int
	)
	err := repo.db.QueryRow(query, ShortURL.URL, ShortURL.Alias).Scan(&id)
	if err != nil {
		return nil, err
	}
	return model.NewShortURL(id, ShortURL.URL, ShortURL.Alias), nil
}
