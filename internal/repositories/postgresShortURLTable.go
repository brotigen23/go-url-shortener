package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

//---------------------- ShortURLs table ----------------------

func (repo PostgresRepository) GetAllShortURL() ([]model.ShortURL, error) {
	ret := []model.ShortURL{}
	query, err := repo.db.Query(`SELECT * FROM Short_URLs`)
	if err != nil {
		return nil, err
	}
	var ID int
	var URL string
	var Alias string
	for query.Next() {
		err = query.Scan(&ID, &URL, &Alias)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *model.NewShortURL(ID, URL, Alias))
	}
	err = query.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

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
	return &model.ShortURL{ID: ID, URL: URL, Alias: Alias}, nil
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
	return &model.ShortURL{ID: ID, URL: URL, Alias: Alias}, nil
}

func (repo PostgresRepository) GetShortURLByURL(URL string) (*model.ShortURL, error) { return nil, nil }

func (repo PostgresRepository) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	var count int
	err := repo.db.QueryRow("SELECT COUNT(*) FROM Short_URLs WHERE URL = $1", ShortURL.URL).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return &ShortURL, fmt.Errorf("URL already exists")
	}
	query := "INSERT INTO Short_URLs(URL, Alias) VALUES($1, $2) RETURNING ID"
	var (
		id int
	)
	err = repo.db.QueryRow(query, ShortURL.URL, ShortURL.Alias).Scan(&id)
	if err != nil {
		return nil, err
	}
	return model.NewShortURL(id, ShortURL.URL, ShortURL.Alias), nil
}
