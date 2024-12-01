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
	repo        repositories.Repo
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
	usersShortURLs, err := service.repository.GetUsersShortURLSByUserID(user.Id)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}
	// Get user's shortURL
	var urls []*model.ShortURL
	for _, userShortURL := range usersShortURLs {
		url, err := service.repository.GetShortURLByID(userShortURL.URL_ID)
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
	_, err = service.repository.SaveUserShortURL(*model.NewUsers_ShortURLs(0, user.Id, shortURL.Id))
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
	usersURLID, err := service.repository.GetUsersShortURLSByUserID(user.Id)
	if err != nil {
		return "", fmt.Errorf("no saved urls")
	}
	// Get user's shortURL
	var urls []model.ShortURL
	for _, urlID := range usersURLID {
		url, err := service.repository.GetShortURLByID(urlID.URL_ID)
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
	usersURLID, err := service.repository.GetUsersShortURLSByUserID(user.Id)
	if err != nil {
		return nil, err
	}
	// Get user's shortURL
	var urls []model.ShortURL
	for _, urlID := range usersURLID {
		url, err := service.repository.GetShortURLByID(urlID.URL_ID)
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

// ---------------------------- DEPRECATED ----------------------------

func (service ServiceShortener) GetURLByAlia(userName string, alias string) (string, error) {
	_, err := service.repository.GetUserByName(userName)
	if err != nil {
		return "", err
	}

	ret, err := service.repository.GetShortURLByAlias(alias)
	if err != nil {
		return "", err
	}
	fmt.Println(ret)
	return ret.URL, nil
}

func (s *ServiceShortener) GetURLByAlias(alias string) (string, error) {
	ret, err := s.repo.GetByAlias(alias)
	if err != nil {
		return "", err
	}
	fmt.Println(ret)
	return ret.URL, nil
}

func (s *ServiceShortener) GetAliasByURL(url string) (string, error) {
	ret, err := s.repo.GetByURL(url)
	if err != nil {
		return "", err
	}
	return ret.Alias, nil
}

func (s *ServiceShortener) Save(userName string, url string) (string, error) {
	model := model.NewShortURL(0, url, utils.NewRandomString(s.lengthAlias))
	//isExists, err := s.repo.GetByURL(url)

	err := s.repo.Save(*model)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "aliases_url_key"` {
		model, _ = s.repo.GetByURL(model.URL)
	}
	return model.Alias, err
}

func (s *ServiceShortener) GetAll() *[]model.ShortURL {
	return s.repo.GetAll()
}

func (s *ServiceShortener) CheckDBConnection() error {
	return s.repo.CheckDBConnection()
}

func (s *ServiceShortener) Close() {
	s.repo.Close()
}

func (s *ServiceShortener) SaveUsersURL(userID string, alias string) error {
	return s.repo.SaveUserURL(userID, alias)
}

func (s *ServiceShortener) GetUserURL(userID string) ([]model.ShortURL, error) {
	return s.repo.GetUserURL(userID)
}
