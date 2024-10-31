package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

// привязать размер алиасов к конструктору репозитория

type inMemoryRepo struct {
	aliases []model.Alias
}

func NewInMemoryRepository() *inMemoryRepo {
	return &inMemoryRepo{
		aliases: []model.Alias{},
	}
}

func (repo inMemoryRepo) GetByAlias(alias string) (*model.Alias, error) {
	for _, v := range repo.aliases {
		if v.GetAlias() == alias {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (repo inMemoryRepo) GetByURL(url string) (*model.Alias, error) {
	for _, v := range repo.aliases {
		if v.GetURL() == url {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (repo *inMemoryRepo) Save(model model.Alias) error {
	for _, v := range repo.aliases {
		if v == model {
			return fmt.Errorf("already exist")
		}
	}
	repo.aliases = append(repo.aliases, model)
	return nil
}
