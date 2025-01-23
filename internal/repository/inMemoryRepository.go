package repository

import (
	"github.com/brotigen23/go-url-shortener/internal/model"
)

type inMemoryRepo struct {
	Users          []model.User
	ShortURLs      []model.ShortURL
	UsersShortURLs []model.UsersShortURLs
}

func NewInMemoryRepo(shortURLs []model.ShortURL, users []model.User, userURLs []model.UsersShortURLs) *inMemoryRepo {
	return &inMemoryRepo{
		ShortURLs:      shortURLs,
		Users:          users,
		UsersShortURLs: userURLs,
	}
}

func NewInMemoryRepository(shortURLs []model.ShortURL) *inMemoryRepo {
	return &inMemoryRepo{
		ShortURLs: shortURLs,
	}
}
func (r *inMemoryRepo) Close() error {
	return nil
}

func (r *inMemoryRepo) CheckDBConnection() error { return nil }
