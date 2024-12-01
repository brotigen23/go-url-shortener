package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

func (repo inMemoryRepo) GetAllUsersShortURLS() ([]model.Users_ShortURLs, error) { return nil, nil }

func (repo inMemoryRepo) GetUsersShortURLSByID(ID int) (*model.Users_ShortURLs, error) {
	return nil, nil
}
func (repo inMemoryRepo) GetUsersShortURLSByUserID(userID int) ([]model.Users_ShortURLs, error) {
	return nil, nil
}
func (repo inMemoryRepo) GetUsersShortURLSByURLID(urlID int) (*model.Users_ShortURLs, error) {
	return nil, nil
}

func (repo inMemoryRepo) SaveUserShortURL(shortURL model.Users_ShortURLs) (*model.Users_ShortURLs, error) {
	return nil, nil
}
