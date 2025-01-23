package repository

import "github.com/brotigen23/go-url-shortener/internal/model"

func (r *inMemoryRepo) GetAllUsersShortURLS() ([]model.UsersShortURLs, error) { return nil, nil }

func (r *inMemoryRepo) GetUsersShortURLSByID(ID int) (*model.UsersShortURLs, error) {
	return nil, nil
}
func (r *inMemoryRepo) GetUsersShortURLSByUserID(userID int) ([]model.UsersShortURLs, error) {
	var ret []model.UsersShortURLs
	for _, UserURL := range r.UsersShortURLs {
		if UserURL.UserID == userID {
			ret = append(ret, UserURL)
		}
	}
	return ret, nil
}
func (r *inMemoryRepo) GetUsersShortURLSByURLID(urlID int) (*model.UsersShortURLs, error) {
	return nil, nil
}

func (r *inMemoryRepo) SaveUserShortURL(shortURL model.UsersShortURLs) (*model.UsersShortURLs, error) {
	shortURL.ID = len(r.UsersShortURLs)
	r.UsersShortURLs = append(r.UsersShortURLs, shortURL)
	return &shortURL, nil
}
