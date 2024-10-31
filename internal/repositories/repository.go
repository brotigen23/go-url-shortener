package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

type Repository interface {
	GetByAlias(alias string) (*model.Alias, error)
	GetByURL(url string) (*model.Alias, error)
	Save(model model.Alias) error
}
