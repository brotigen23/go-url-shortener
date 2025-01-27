package repository

import "github.com/brotigen23/go-url-shortener/internal/model"

type Repository interface {
	Create(shortURL model.ShortURL) error
	GetAll() ([]model.ShortURL, error)
	GetByUser(username string) ([]model.ShortURL, error)
	GetByURL(url string) (*model.ShortURL, error)
	GetByAlias(alias string) (*model.ShortURL, error)

	Update(username string, shortURL model.ShortURL) error

	Delete(username string, shortURL []model.ShortURL) error
}
