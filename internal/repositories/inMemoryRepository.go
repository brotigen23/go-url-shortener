package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

type inMemoryRepo struct {
	Users           []model.User
	ShortURLs       []model.ShortURL
	Users_ShortURLs []model.Users_ShortURLs
	aliases         []model.ShortURL
}

func NewInMemoryRepo(shortURLs []model.ShortURL, users []model.User, userURLs []model.Users_ShortURLs) *inMemoryRepo {
	return &inMemoryRepo{
		ShortURLs: shortURLs,
		Users: users,
		Users_ShortURLs: userURLs,
	}
}

func NewInMemoryRepository(shortURLs []model.ShortURL) *inMemoryRepo {
	return &inMemoryRepo{
		ShortURLs: shortURLs,
	}
}
func (repo inMemoryRepo) GetByAlias(alias string) (*model.ShortURL, error) {
	for _, v := range repo.aliases {
		if v.GetAlias() == alias {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (repo inMemoryRepo) GetByURL(url string) (*model.ShortURL, error) {
	for _, v := range repo.aliases {
		if v.GetURL() == url {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (repo *inMemoryRepo) Save(model model.ShortURL) error {
	for _, v := range repo.aliases {
		if v.GetURL() == model.GetURL() {
			return fmt.Errorf("already exist")
		}
	}
	repo.aliases = append(repo.aliases, model)
	return nil
}
func (repo *inMemoryRepo) GetAll() *[]model.ShortURL {
	return &repo.aliases
}

func (repo *inMemoryRepo) Migrate(aliases []model.ShortURL) {
	repo.aliases = append(repo.aliases, aliases...)
}

func (repo *inMemoryRepo) Close() error {
	return nil
}

func (repo *inMemoryRepo) CheckDBConnection() error { return nil }

func (repo *inMemoryRepo) SaveUser1(id string) error { return nil }

func (repo *inMemoryRepo) GetUserByID1(userID string) error { return nil }

func (repo *inMemoryRepo) SaveUserURL1(userID string, alias string) error {
	return nil
}
func (repo *inMemoryRepo) GetUserURL1(userID string) ([]model.ShortURL, error) {
	return nil, nil
}
