package service

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"go.uber.org/zap"
)

// Сервис бизнес логики
type Service struct {
	repository  repository.Repository
	lengthAlias int
	config      *config.Config
	logger      *zap.SugaredLogger
}

// Конструктор Service
func New(config *config.Config, logger *zap.SugaredLogger, repository repository.Repository) *Service {
	return &Service{
		repository:  repository,
		lengthAlias: 8,
		config:      config,
		logger:      logger,
	}
}

// Обновляет длину создаваемых алиасов
func (s *Service) SetLengthAlias(lengthAlias int) {
	s.lengthAlias = lengthAlias
}

// Create short url
func (s Service) CreateShortURL(username string, URL string) (string, error) {
	alias := utils.NewRandomString(s.lengthAlias)
	err := s.repository.Create(model.ShortURL{URL: URL, ShortURL: alias, Username: username})
	if err != nil {
		if err == repository.ErrShortURLAlreadyExists {
			ret, er := s.repository.GetByURL(URL)
			if er != nil {
				return "", er
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

// Возвращает строку соединения с базой данных
func (s Service) GetDSN() string {
	return s.config.DatabaseDSN
}
