package postgres

import (
	"database/sql"

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

func (r repository) Create(shortURL model.ShortURL) error { return nil }

func (r *repository) GetAll() ([]model.ShortURL, error) { return nil, nil }

func (r repository) GetByUser(username string) ([]model.ShortURL, error) { return nil, nil }
func (r repository) GetByURL(url string) (*model.ShortURL, error)        { return nil, nil }
func (r repository) GetByAlias(alias string) (*model.ShortURL, error)    { return nil, nil }

func (r repository) Update(username string, shortURL model.ShortURL) error { return nil }

func (r repository) Delete(username string, shortURL []model.ShortURL) error { return nil }
