package postgres

import (
	"database/sql"
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

type repository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func New(db *sql.DB, logger *zap.SugaredLogger) *repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r repository) Create(shortURL model.ShortURL) error {
	query := `	
	INSERT INTO short_url(url, short_url, username) 
	VALUES($1, $2, $3)`

	_, err := r.db.Exec(query, shortURL.URL, shortURL.ShortURL, shortURL.Username)
	if err != nil {
		r.logger.Errorln(err)
		return err
	}
	return nil
}

func (r repository) GetAll() ([]model.ShortURL, error) {
	ret := []model.ShortURL{}
	query := `
	SELECT id, url, short_url, username, is_deleted 
	FROM short_url`

	row, err := r.db.Query(query)
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	var ID int
	var URL, short_url, username string
	var IsDeleted bool

	for row.Next() {
		err := row.Scan(&ID, &URL, &short_url, &username, &IsDeleted)
		if err != nil {
			r.logger.Errorln(err)
			return nil, err
		}
		ret = append(ret, model.ShortURL{ID: ID, URL: URL, ShortURL: short_url, Username: username, IsDeleted: IsDeleted})
	}
	err = row.Err()
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	return ret, nil
}

func (r repository) GetByUser(username string) ([]model.ShortURL, error) {
	ret := []model.ShortURL{}
	query := `
	SELECT id, url, short_url, is_deleted 
	FROM short_url
	WHERE username = $1`

	row, err := r.db.Query(query, username)
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	var ID int
	var URL, shortURL string
	var IsDeleted bool

	for row.Next() {
		err := row.Scan(&ID, &URL, &shortURL, &IsDeleted)
		if err != nil {
			r.logger.Errorln(err)
			return nil, err
		}
		ret = append(ret, model.ShortURL{ID: ID, URL: URL, ShortURL: shortURL, Username: username, IsDeleted: IsDeleted})
	}
	err = row.Err()
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	return ret, nil
}
func (r repository) GetByURL(url string) (*model.ShortURL, error) {
	query := `
	SELECT id, short_url, username, is_deleted 
	FROM short_url
	WHERE url = $1`

	row := r.db.QueryRow(query, url)
	var ID int
	var shortURL, username string
	var IsDeleted bool

	err := row.Scan(&ID, &shortURL, &username, &IsDeleted)
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	err = row.Err()
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}

	return &model.ShortURL{ID: ID, URL: url, ShortURL: shortURL, Username: username, IsDeleted: IsDeleted}, nil
}
func (r repository) GetByAlias(alias string) (*model.ShortURL, error) {
	query := `
	SELECT id, url, username, is_deleted 
	FROM short_url
	WHERE short_url = $1`

	row := r.db.QueryRow(query, alias)
	var ID int
	var url, username string
	var IsDeleted bool

	err := row.Scan(&ID, &url, &username, &IsDeleted)
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	err = row.Err()
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}

	return &model.ShortURL{ID: ID, URL: url, ShortURL: alias, Username: username, IsDeleted: IsDeleted}, nil
}

func (r repository) Update(username string, shortURL model.ShortURL) error { return nil }

func (r repository) Delete(username string, shortURL []model.ShortURL) error {
	query := `
	UPDATE short_url 
	SET is_deleted = true
	WHERE short_url = $1 AND username = $2`

	_, err := r.db.Exec(query, fmt.Sprint(shortURL), username)
	if err != nil {
		r.logger.Errorln(err)
		return err
	}
	return nil

}
