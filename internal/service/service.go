package service

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"go.uber.org/zap"
)

type Service struct {
	repository  repository.Repository
	lengthAlias int
	config      *config.Config
	logger      *zap.SugaredLogger
}

func New(config *config.Config, logger *zap.SugaredLogger, repository repository.Repository) (*Service, error) {
	return &Service{
		repository:  repository,
		lengthAlias: 8,
		config:      config,
		logger:      logger,
	}, nil
}

func (s *Service) SetLengthAlias(lengthAlias int) {
	s.lengthAlias = lengthAlias
}

// Create short url
func (s Service) CreateShortURL(username string, URL string) (string, error) {
	alias := utils.NewRandomString(s.lengthAlias)
	err := s.repository.Create(model.ShortURL{URL: URL, ShortURL: alias, Username: username})
	if err != nil {
		if err == repository.ErrShortURLAlreadyExists {
			ret, err := s.repository.GetByURL(URL)
			if err != nil {
				return "", err
			}
			return ret.ShortURL, ErrShortURLAlreadyExists
		}
		return "", err
	}
	s.logger.Infoln(URL, "saved", "as", alias, "by", username)
	return alias, nil
}

// Create short urls
func (s Service) CreateShortURLs(username string, URLs []string) (map[string]string, error) {
	var conflict error
	ret := make(map[string]string)
	for _, url := range URLs {
		shortURL, err := s.CreateShortURL(username, url)
		if err != nil {
			if err == ErrShortURLAlreadyExists {
				conflict = ErrShortURLAlreadyExists
			} else {
				return nil, err
			}
		}
		ret[url] = shortURL
	}
	return ret, conflict
}

// Return URL founded by alias
func (s Service) GetShortURL(alias string) (string, error) {
	shortURL, err := s.repository.GetByAlias(alias)
	if err != nil {
		if err == repository.ErrNoFound {
			return "", ErrShortURLNotFound
		}
		s.logger.Errorln(err)
		return "", err
	}
	return shortURL.URL, nil
}

// Return all URLs saved by user
func (s Service) GetShortURLs(username string) (map[string]string, error) {
	s.logger.Debugln("get", username, "'s urls")
	ret := make(map[string]string)
	shortURLs, err := s.repository.GetByUser(username)
	if err != nil {
		if err == repository.ErrNoFound {
			return nil, ErrShortURLNotFound
		}
		s.logger.Errorln(err)
		return nil, err
	}
	for _, v := range shortURLs {
		ret[v.URL] = v.ShortURL
	}
	return ret, nil
}

// Delete short urls saved by user
func (s Service) DeleteShortURLs(username string, aliases []string) error {
	shortURLs := model.NewShortURLs(aliases)
	err := s.repository.Delete(username, shortURLs)
	if err == repository.ErrNoFound {
		return ErrShortURLNotFound
	}
	return err
}

// Check if short url is deleted
func (s Service) IsShortURLDeleted(alias string) (bool, error) {
	shortURL, err := s.repository.GetByAlias(alias)
	if err != nil {
		if err == repository.ErrNoFound {
			return false, ErrShortURLNotFound
		}
		return false, err
	}
	return shortURL.IsDeleted, nil
}

func (s Service) GetDSN() string {
	return s.config.DatabaseDSN
}

// -----------------------------------------------------------------
//
//	DEPRECATED
//
// -----------------------------------------------------------------
/*
func SaveURL(userName string, URL string) (string, error) {
	user, err := s.repository.GetUserByName(userName)
	if err != nil {
		return "", err
	}
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
	_, err = s.repository.SaveUserShortURL(*model.NewUsersShortURLs(0, user.ID, shortURL.ID))
	if err != nil {
		return "", nil
	}
	s.logger.Infoln(shortURL.URL, "saved", URL, "as", alias, "by", userName)
	return shortURL.Alias, nil
}

func SaveURLs(userName string, URLs []string) (map[string]string, error) {
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

func GetURL(alias string) (string, error) {
	ret, err := s.repository.GetShortURLByAlias(alias)
	if err != nil {
		return "", err
	}
	return ret.URL, nil
}

func GetURLs(userName string) (map[string]string, error) {
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
	var urls []model.ShortURL1
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

func DeleteURLs(userName string, aliases []string) error {
	err := s.repository.DeleteShortURLByAliases(aliases)
	if err != nil {
		return err
	}
	return nil
}

func IsDeleted(alias string) (bool, error) {
	d, err := s.repository.GetShortURLByAlias(alias)
	if err != nil {
		return false, err
	}
	fmt.Println(d)
	return d.IsDeleted, nil
}

func CheckDBConnection() error {
	return nil
}

func GetBaseURL() string {
	return s.config.BaseURL
}
*/
