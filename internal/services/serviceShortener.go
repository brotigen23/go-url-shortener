package services

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repositories"
	"github.com/brotigen23/go-url-shortener/internal/utils"
)

type ServiceShortener struct {
	repo        repositories.Repository
	lengthAlias int
}

func NewService(config *config.Config, lengthAlias int, a []model.Alias) (*ServiceShortener, error) {
	if config.DatabaseDSN != "" {
		db, err := repositories.NewPostgresRepository(config.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		return &ServiceShortener{
			repo:        db,
			lengthAlias: lengthAlias,
		}, nil
	} else {
		return &ServiceShortener{
			repo:        repositories.NewInMemoryRepository(a),
			lengthAlias: lengthAlias,
		}, nil
	}
}

func (s *ServiceShortener) GetURLByAlias(alias string) (string, error) {
	ret, err := s.repo.GetByAlias(alias)
	if err != nil {
		return "", err
	}
	fmt.Println(ret)
	return ret.URL, nil
}

func (s *ServiceShortener) GetAliasByURL(url string) (string, error) {
	ret, err := s.repo.GetByURL(url)
	if err != nil {
		return "", err
	}
	return ret.Alias, nil
}

func (s *ServiceShortener) Save(url string) (string, error) {
	model := model.NewAlias(url, utils.NewRandomString(s.lengthAlias))

	err := s.repo.Save(*model)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "aliases_url_key"` {
		model, _ = s.repo.GetByURL(model.URL)
	}
	return model.Alias, err
}

func (s *ServiceShortener) GetAll() *[]model.Alias {
	return s.repo.GetAll()
}

func (s *ServiceShortener) CheckDBConnection() error {
	return s.repo.CheckDBConnection()
}

func (s *ServiceShortener) Close() {
	s.repo.Close()
}
