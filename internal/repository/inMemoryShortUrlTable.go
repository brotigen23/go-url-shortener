package repository

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

func (r *inMemoryRepo) GetAllShortURL() ([]model.ShortURL, error) {
	return r.ShortURLs, nil
}

func (r *inMemoryRepo) GetShortURLByID(id int) (*model.ShortURL, error) {
	for _, shortURL := range r.ShortURLs {
		if shortURL.ID == id {
			return &shortURL, nil
		}
	}
	return nil, fmt.Errorf("URL not found")
}
func (r *inMemoryRepo) GetShortURLByAlias(alias string) (*model.ShortURL, error) {
	for _, shortURL := range r.ShortURLs {
		if shortURL.Alias == alias {
			return &shortURL, nil
		}
	}
	return nil, fmt.Errorf("URL not found")
}
func (r *inMemoryRepo) GetShortURLByURL(URL string) (*model.ShortURL, error) {
	for _, shortURL := range r.ShortURLs {
		if shortURL.URL == URL {
			return &shortURL, nil
		}
	}
	return nil, fmt.Errorf("URL not found")
}

func (r *inMemoryRepo) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	fmt.Println(r.ShortURLs)
	for _, url := range r.ShortURLs {
		if url.URL == ShortURL.URL {
			return &url, fmt.Errorf("URL already exists")
		}
	}

	shortURL := model.NewShortURL(len(r.ShortURLs), ShortURL.URL, ShortURL.Alias)
	r.ShortURLs = append(r.ShortURLs, *shortURL)
	return shortURL, nil
}

func (r *inMemoryRepo) DeleteShortURLByAlias(Alias string) error       { return nil }
func (r *inMemoryRepo) DeleteShortURLByAliases(Aliases []string) error { return nil }
