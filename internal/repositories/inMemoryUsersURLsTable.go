package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

func (repo *inMemoryRepo) GetAllUsersShortURLS() ([]model.UsersShortURLs, error) { return nil, nil }

func (repo *inMemoryRepo) GetUsersShortURLSByID(ID int) (*model.UsersShortURLs, error) {
	return nil, nil
}
func (repo *inMemoryRepo) GetUsersShortURLSByUserID(userID int) ([]model.UsersShortURLs, error) {
	var ret []model.UsersShortURLs
	for _, UserURL := range repo.UsersShortURLs {
		if UserURL.UserID == userID {
			ret = append(ret, UserURL)
		}
	}
	return ret, nil
}
func (repo *inMemoryRepo) GetUsersShortURLSByURLID(urlID int) (*model.UsersShortURLs, error) {
	return nil, nil
}

func (repo *inMemoryRepo) SaveUserShortURL(shortURL model.UsersShortURLs) (*model.UsersShortURLs, error) {
	shortURL.ID = len(repo.UsersShortURLs)
	repo.UsersShortURLs = append(repo.UsersShortURLs, shortURL)
	return &shortURL, nil
}
