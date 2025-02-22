package memory

import (
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repository"
)

// Repository с реализацией сохранения данных в память
type Repository struct {
	shortURLs []model.ShortURL
}

// Конструктор Repository
func New(shortURLs []model.ShortURL) *Repository {
	return &Repository{
		shortURLs: shortURLs,
	}
}

// Возвращает все имеющиеся данные
func (r Repository) GetAll() ([]model.ShortURL, error) {
	return r.shortURLs, nil
}

// Создает новую ссылку
func (r *Repository) Create(shortURL model.ShortURL) error {
	for _, v := range r.shortURLs {
		if v.URL == shortURL.URL {
			return repository.ErrShortURLAlreadyExists
		}
	}
	id := len(r.shortURLs)
	shortURL.ID = id
	r.shortURLs = append(r.shortURLs, shortURL)
	return nil
}

// Возвращает все ссылки, сохраненные определенным пользователем
func (r Repository) GetByUser(username string) ([]model.ShortURL, error) {
	var ret []model.ShortURL
	for _, v := range r.shortURLs {
		if v.Username == username {
			ret = append(ret, v)
		}
	}
	if len(ret) == 0 {
		return nil, repository.ErrNoFound
	}
	return ret, nil
}

// Возвращает сущность ссылки по входящему URL
func (r Repository) GetByURL(url string) (*model.ShortURL, error) {
	for _, v := range r.shortURLs {
		if v.URL == url {
			return &v, nil
		}
	}
	return nil, repository.ErrNoFound
}

// Возвращает сущность ссылки по входящему Alias
func (r Repository) GetByAlias(alias string) (*model.ShortURL, error) {
	for _, v := range r.shortURLs {
		if v.ShortURL == alias {
			return &v, nil
		}
	}
	return nil, repository.ErrNoFound
}

// Обновляет входящую сущность
func (r *Repository) Update(username string, shortURL model.ShortURL) error {
	for i, v := range r.shortURLs {
		if v == shortURL {
			r.shortURLs[i] = shortURL
			return nil
		}
	}
	return repository.ErrNoFound
}

// Удаляет входящую сущность
func (r *Repository) Delete(username string, shortURL []model.ShortURL) error {
	for _, v := range shortURL {
		for i, k := range r.shortURLs {
			if k == v {
				r.shortURLs[i].IsDeleted = true
			}
		}
	}
	return nil
}
