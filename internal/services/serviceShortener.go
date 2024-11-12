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

func NewService(config *config.Config, lengthAlias int, a []model.Alias) *ServiceShortener {
	if config.DatabaseDSN != "" {
		fmt.Println("Create DB")
		return &ServiceShortener{
			repo:        repositories.NewPostgresRepository(config.DatabaseDSN),
			lengthAlias: lengthAlias,
		}
	} else {
		fmt.Println("Set inMem")
		return &ServiceShortener{
			repo:        repositories.NewInMemoryRepository(a),
			lengthAlias: lengthAlias,
		}
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
