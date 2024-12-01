package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

func (repo *inMemoryRepo) GetAllShortURL() ([]model.ShortURL, error) { return nil, nil }

func (repo *inMemoryRepo) GetShortURLByID(id int) (*model.ShortURL, error) {
	for _, shortURL := range repo.ShortURLs {
		if shortURL.ID == id {
			return &shortURL, nil
		}
	}
	return nil, fmt.Errorf("URL not found")
}
func (repo *inMemoryRepo) GetShortURLByAlias(alias string) (*model.ShortURL, error) {
	for _, shortURL := range repo.ShortURLs {
		if shortURL.Alias == alias {
			return &shortURL, nil
		}
	}
	return nil, fmt.Errorf("URL not found")
}
func (repo *inMemoryRepo) GetShortURLByURL(URL string) (*model.ShortURL, error) {
	for _, shortURL := range repo.ShortURLs {
		if shortURL.URL == URL {
			return &shortURL, nil
		}
	}
	return nil, fmt.Errorf("URL not found")
}

func (repo *inMemoryRepo) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	fmt.Println(repo.ShortURLs)
	for _, url := range repo.ShortURLs {
		if url.URL == ShortURL.URL {
			return &url, fmt.Errorf("URL already exists")
		}
	}

	shortURL := model.NewShortURL(len(repo.ShortURLs), ShortURL.URL, ShortURL.Alias)
	repo.ShortURLs = append(repo.ShortURLs, *shortURL)
	return shortURL, nil
}
