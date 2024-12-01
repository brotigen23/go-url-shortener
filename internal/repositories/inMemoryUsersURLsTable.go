package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

func (repo inMemoryRepo) GetAllUsersShortURLS() ([]model.UsersShortURLs, error) { return nil, nil }

func (repo inMemoryRepo) GetUsersShortURLSByID(ID int) (*model.UsersShortURLs, error) {
	return nil, nil
}
func (repo inMemoryRepo) GetUsersShortURLSByUserID(userID int) ([]model.UsersShortURLs, error) {
	return nil, nil
}
func (repo inMemoryRepo) GetUsersShortURLSByURLID(urlID int) (*model.UsersShortURLs, error) {
	return nil, nil
}

func (repo inMemoryRepo) SaveUserShortURL(shortURL model.UsersShortURLs) (*model.UsersShortURLs, error) {
	return nil, nil
}
