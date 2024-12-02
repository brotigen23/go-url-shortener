package repositories

import (
	"fmt"
	"strings"

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
	var IsDeleted bool

	for query.Next() {
		err := query.Scan(&ID, &URL, &Alias, &IsDeleted)
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
	var IsDeleted bool
	err := query.Scan(&ID, &URL, &Alias, &IsDeleted)
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
	var IsDeleted bool
	err := query.Scan(&ID, &URL, &Alias, &IsDeleted)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.ShortURL{ID: ID, URL: URL, Alias: Alias, IsDeleted: IsDeleted}, nil
}

func (repo PostgresRepository) GetShortURLByURL(url string) (*model.ShortURL, error) {
	query := repo.db.QueryRow(`SELECT * FROM Short_URLs WHERE URL = $1`, url)
	var ID int
	var URL string
	var Alias string
	var IsDeleted bool
	err := query.Scan(&ID, &URL, &Alias, &IsDeleted)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &model.ShortURL{ID: ID, URL: URL, Alias: Alias}, nil

}

func (repo PostgresRepository) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	query := "INSERT INTO Short_URLs(URL, Alias) VALUES($1, $2) RETURNING ID"
	var (
		id int
	)
	err := repo.db.QueryRow(query, ShortURL.URL, ShortURL.Alias).Scan(&id)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "short_urls_url_key"` {
			ret, e := repo.GetShortURLByURL(ShortURL.URL)
			if e != nil {
				return nil, e
			}
			return ret, err
		}
		fmt.Println(err.Error())
		return nil, err
	}
	return model.NewShortURL(id, ShortURL.URL, ShortURL.Alias), nil
}

func (repo PostgresRepository) DeleteShortURLByAlias(Alias string) error { return nil }

func (repo PostgresRepository) DeleteShortURLByAliases(Aliases []string) error {
	aliases := strings.Join(Aliases[:], ",")
	query := "UPDATE Short_URLs SET Is_Deleted = TRUE WHERE Alias IN ($1) "
	err := repo.db.QueryRow(query, aliases).Err()
	if err != nil {
		return err
	}
	return err
}
