package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

type Repository interface {
	GetByAlias(alias string) (*model.Alias, error)
	GetByURL(url string) (*model.Alias, error)
	GetAll() *[]model.Alias
	Save(model model.Alias) error
	Migrate(model []model.Alias)
	CheckDBConnection() error
	SaveUser(id string) error
	GetUserByID(userID string) error
	SaveUserURL(userID string, alias string) error
	GetUserURL(userID string) ([]model.Alias, error)

	Close()
}
