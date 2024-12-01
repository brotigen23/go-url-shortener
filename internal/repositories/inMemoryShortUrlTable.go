package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

func (repo *inMemoryRepo) GetAllShortURL() ([]model.ShortURL, error) { return nil, nil }

func (repo *inMemoryRepo) GetShortURLByID(id int) (*model.ShortURL, error)          { return nil, nil }
func (repo *inMemoryRepo) GetShortURLByAlias(alias string) (*model.ShortURL, error) { return nil, nil }
func (repo *inMemoryRepo) GetShortURLByURL(URL string) (*model.ShortURL, error)     { return nil, nil }

func (repo *inMemoryRepo) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	shortURL := model.NewShortURL(len(repo.ShortURLs), ShortURL.URL, ShortURL.Alias)
	repo.ShortURLs = append(repo.ShortURLs, *shortURL)
	return shortURL, nil
}
