package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"go.uber.org/zap"
)

// Repository с реализацией сохранения данных в базу данных Postgres
type Repository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

// Конструктор Repository
func New(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Создает новую ссылку
func (r *Repository) Create(shortURL model.ShortURL) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `	
	INSERT INTO short_url(url, short_url, username) 
	VALUES($1, $2, $3)`

	_, err = tx.Exec(query, shortURL.URL, shortURL.ShortURL, shortURL.Username)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "short_url_url_key"` {
			err = repository.ErrShortURLAlreadyExists
		}
		e := tx.Rollback()
		if e != nil {
			return e
		}
		return err
	}
	err = tx.Commit()
	return err
}

// Возвращает все имеющиеся данные
func (r *Repository) GetAll() ([]model.ShortURL, error) {
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
	var URL, shortURL, username string
	var IsDeleted bool

	for row.Next() {
		err = row.Scan(&ID, &URL, &shortURL, &username, &IsDeleted)
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

// Возвращает все ссылки, сохраненные определенным пользователем
func (r *Repository) GetByUser(username string) ([]model.ShortURL, error) {
	ret := make([]model.ShortURL, 0, 100)
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
		err = row.Scan(&ID, &URL, &shortURL, &IsDeleted)
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

// Возвращает сущность ссылки по входящему URL
func (r *Repository) GetByURL(url string) (*model.ShortURL, error) {
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

// Возвращает сущность ссылки по входящему Alias
func (r *Repository) GetByAlias(alias string) (*model.ShortURL, error) {
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

// Обновляет входящую сущность
func (r *Repository) Update(username string, shortURL model.ShortURL) error { return nil }

// Удаляет входящую сущность
func (r *Repository) Delete(username string, shortURL []model.ShortURL) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// ------------------------------------------------
	// generate query
	// ------------------------------------------------
	aliases := make([]string, 0)
	for i := range shortURL {
		aliases = append(aliases, "")
		aliases[i] = `'` + shortURL[i].ShortURL + `'`
	}
	toDelete := strings.Join(aliases[:], ",")
	r.logger.Debugln("url to delete", toDelete)
	query := fmt.Sprintf(`
	UPDATE short_url 
	SET is_deleted = TRUE
	WHERE short_url IN (%s)`, toDelete)

	// ------------------------------------------------
	// do query
	// ------------------------------------------------
	_, err = tx.Exec(query)

	if err != nil {
		r.logger.Errorln(err)
		e := tx.Rollback()
		if e != nil {
			return e
		}
		return err
	}
	err = tx.Commit()
	return err
}
