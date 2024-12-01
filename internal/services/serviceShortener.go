package services

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repositories"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"go.uber.org/zap"
)

type ServiceShortener struct {
	repository  repositories.Repository
	lengthAlias int
}

func NewService(config *config.Config, lengthAlias int, a []model.ShortURL, logger *zap.Logger, repository repositories.Repository) (*ServiceShortener, error) {
	return &ServiceShortener{
		repository:  repository,
		lengthAlias: lengthAlias,
	}, nil
}

func (service ServiceShortener) SaveURL(userName string, URL string) (string, error) {
	// Get user entity
	user, err := service.repository.GetUserByName(userName)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}
	// Get user's URL IDs
	usersShortURLs, err := service.repository.GetUsersShortURLSByUserID(user.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}
	// Get user's shortURL
	var urls []*model.ShortURL
	for _, userShortURL := range usersShortURLs {
		url, err := service.repository.GetShortURLByID(userShortURL.URLID)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return "", nil
		}
		urls = append(urls, url)
	}
	// Check if URL already exists
	for _, url := range urls {
		if url.URL == URL {
			return url.Alias, fmt.Errorf("URL already exists")
		}
	}
	// Create new shortURL
	alias := utils.NewRandomString(service.lengthAlias)
	shortURL, err := service.repository.SaveShortURL(*model.NewShortURL(0, URL, alias))
	if err != nil {
		return "", nil
	}
	// Create relation User <-> URL
	_, err = service.repository.SaveUserShortURL(*model.NewUsersShortURLs(0, user.ID, shortURL.ID))
	if err != nil {
		return "", nil
	}
	return shortURL.Alias, nil
}

// For BATCH
func (service ServiceShortener) SaveURLs(userName string, URLs []string) (map[string]string, error) {
	ret := make(map[string]string)
	for _, url := range URLs {
		shortURL, err := service.SaveURL(userName, url)
		if err != nil {
			return nil, err
		}
		ret[url] = shortURL
	}
	return ret, nil
}

func (service ServiceShortener) GetURL(userName string, alias string) (string, error) {
	// Get user entity
	user, err := service.repository.GetUserByName(userName)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	// Get user's URL IDs
	usersURLID, err := service.repository.GetUsersShortURLSByUserID(user.ID)
	if err != nil {
		return "", fmt.Errorf("no saved urls")
	}
	// Get user's shortURL
	var urls []model.ShortURL
	for _, urlID := range usersURLID {
		url, err := service.repository.GetShortURLByID(urlID.URLID)
		if err != nil {
			return "", err
		}
		urls = append(urls, *url)
	}
	// Search url with alias
	for _, url := range urls {
		if url.Alias == alias {
			return url.URL, nil
		}
	}
	return "", fmt.Errorf("url not found")
}

func (service ServiceShortener) GetURLs(userName string) (map[string]string, error) {
	ret := make(map[string]string)

	// Get user entity
	user, err := service.repository.GetUserByName(userName)
	if err != nil {
		return nil, err
	}
	// Get user's URL IDs
	usersURLID, err := service.repository.GetUsersShortURLSByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	// Get user's shortURL
	var urls []model.ShortURL
	for _, urlID := range usersURLID {
		url, err := service.repository.GetShortURLByID(urlID.URLID)
		if err != nil {
			return nil, err
		}
		urls = append(urls, *url)
	}
	// Search url with alias
	for _, url := range urls {
		ret[url.URL] = url.Alias
	}
	return ret, nil
}

func (service ServiceShortener) CheckDBConnection() error {
	return service.repository.CheckDBConnection()
}
