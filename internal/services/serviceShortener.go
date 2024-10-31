package services

import (
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repositories"
	"github.com/brotigen23/go-url-shortener/internal/utils"
)

type ServiceShortener struct {
	repo        repositories.Repository
	lengthAlias int
}

func NewService(lengthAlias int, a []model.Alias) *ServiceShortener {
	return &ServiceShortener{
		repo:        repositories.NewInMemoryRepository(a),
		lengthAlias: lengthAlias,
	}
}

func NewMockService(lengthAlias int) *ServiceShortener {
	return &ServiceShortener{
		repo:        repositories.NewMockInMemoryRepository(),
		lengthAlias: lengthAlias,
	}
}

func (s ServiceShortener) GetURLByAlias(alias string) (*model.Alias, error) {
	return s.repo.GetByAlias(alias)
}

func (s ServiceShortener) GetAliasByURL(url string) (*model.Alias, error) {
	return s.repo.GetByURL(url)
}

func (s ServiceShortener) Save(url string) (*model.Alias, error) {
	model := model.NewAlias(url, utils.NewRandomString(s.lengthAlias))

	err := s.repo.Save(*model)

	return model, err
}

func (s ServiceShortener) GetAll() *[]model.Alias {
	return s.repo.GetAll()
}
