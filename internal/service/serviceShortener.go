package service

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"go.uber.org/zap"
)

type ServiceShortener struct {
	repository  repository.Repository
	lengthAlias int
	config      *config.Config
	logger      *zap.SugaredLogger
}

func NewServiceShortener(config *config.Config, lengthAlias int, logger *zap.SugaredLogger, repository repository.Repository) (*ServiceShortener, error) {
	return &ServiceShortener{
		repository:  repository,
		lengthAlias: lengthAlias,
		config:      config,
		logger:      logger,
	}, nil
}

func (s ServiceShortener) SaveURL(userName string, URL string) (string, error) {
	// Get user entity
	user, err := s.repository.GetUserByName(userName)
	if err != nil {
		return "", err
	}
	// Create new shortURL
	alias := utils.NewRandomString(s.lengthAlias)
	shortURL, err := s.repository.SaveShortURL(*model.NewShortURL(0, URL, alias))
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "short_urls_url_key"` {
			s.logger.Debugln("URL:", URL, "already saved as", alias)
			return shortURL.Alias, err
		} else {
			return "", nil
		}
	}
	// Create relation User <-> URL
	_, err = s.repository.SaveUserShortURL(*model.NewUsersShortURLs(0, user.ID, shortURL.ID))
	if err != nil {
		return "", nil
	}
	s.logger.Infoln(userName, "saved", URL, "as", alias)
	return shortURL.Alias, nil
}

// For BATCH
func (s ServiceShortener) SaveURLs(userName string, URLs []string) (map[string]string, error) {
	ret := make(map[string]string)
	for _, url := range URLs {
		shortURL, err := s.SaveURL(userName, url)
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

func (s ServiceShortener) GetURL(alias string) (string, error) {
	ret, err := s.repository.GetShortURLByAlias(alias)
	if err != nil {
		return "", err
	}
	return ret.URL, nil
}

func (s ServiceShortener) GetURLs(userName string) (map[string]string, error) {
	ret := make(map[string]string)

	// Get user entity
	user, err := s.repository.GetUserByName(userName)
	if err != nil {
		return nil, err
	}
	// Get user's URL IDs
	usersURLID, err := s.repository.GetUsersShortURLSByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	// Get user's shortURL
	var urls []model.ShortURL
	for _, urlID := range usersURLID {
		url, err := s.repository.GetShortURLByID(urlID.URLID)
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

func (s ServiceShortener) DeleteURLs(userName string, aliases []string) error {
	err := s.repository.DeleteShortURLByAliases(aliases)
	if err != nil {
		return err
	}
	return nil
}

func (s ServiceShortener) IsDeleted(alias string) (bool, error) {
	d, err := s.repository.GetShortURLByAlias(alias)
	if err != nil {
		return false, err
	}
	fmt.Println(d)
	return d.IsDeleted, nil
}

func (s ServiceShortener) CheckDBConnection() error {
	return s.repository.CheckDBConnection()
}

func (s ServiceShortener) AllURLs() []model.ShortURL {
	ret, err := s.repository.GetAllShortURL()
	if err != nil {
		return nil
	}
	return ret
}

func (s ServiceShortener) GetBaseURL() string {
	return s.config.BaseURL
}
