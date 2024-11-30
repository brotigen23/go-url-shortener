package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

// привязать размер алиасов к конструктору репозитория

type inMemoryRepo struct {
	aliases []model.Alias
}

func NewInMemoryRepository(a []model.Alias) *inMemoryRepo {
	return &inMemoryRepo{
		aliases: a,
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
		if v.GetURL() == model.GetURL() {
			return fmt.Errorf("already exist")
		}
	}
	repo.aliases = append(repo.aliases, model)
	return nil
}
func (repo *inMemoryRepo) GetAll() *[]model.Alias {
	return &repo.aliases
}

func (repo *inMemoryRepo) Migrate(aliases []model.Alias) {
	repo.aliases = append(repo.aliases, aliases...)
}

func (repo *inMemoryRepo) Close() {
}

func (repo *inMemoryRepo) CheckDBConnection() error { return nil }

func (repo *inMemoryRepo) SaveUser(id string) error { return nil }

func (repo *inMemoryRepo) GetUserByID(userID string) error { return nil }

func (repo *inMemoryRepo) SaveUserURL(userID string, alias string) error {
	return nil
}
func (repo *inMemoryRepo) GetUserURL(userID string) ([]model.Alias, error) {
	return nil, nil
}
