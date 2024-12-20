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
	for _, url := range a {
		_, err := repository.SaveShortURL(url)
		if err != nil {
			return nil, err
		}
	}
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
	// Create new shortURL
	alias := utils.NewRandomString(service.lengthAlias)
	shortURL, err := service.repository.SaveShortURL(*model.NewShortURL(0, URL, alias))
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "short_urls_url_key"` {
			return shortURL.Alias, err
		} else {
			return "", nil
		}
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
			if err.Error() == "URL already exists" {
				return nil, err
			} else {
				return nil, err
			}

		}
		ret[url] = shortURL
	}
	return ret, nil
}

func (service ServiceShortener) GetURL(alias string) (string, error) {
	ret, err := service.repository.GetShortURLByAlias(alias)
	if err != nil {
		return "", err
	}
	return ret.URL, nil
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

func (service ServiceShortener) DeleteURLs(userName string, aliases []string) error {
	err := service.repository.DeleteShortURLByAliases(aliases)
	if err != nil {
		return err
	}
	return nil
}

func (service ServiceShortener) IsDeleted(alias string) (bool, error) {
	d, err := service.repository.GetShortURLByAlias(alias)
	if err != nil {
		return false, err
	}
	fmt.Println(d)
	return d.IsDeleted, nil
}

func (service ServiceShortener) CheckDBConnection() error {
	return service.repository.CheckDBConnection()
}

func (service ServiceShortener) AllURLs() []model.ShortURL {
	ret, err := service.repository.GetAllShortURL()
	if err != nil {
		return nil
	}
	return ret
}
