package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

type MockRepository struct {
	userTable      []model.User
	shortURLTable  []model.ShortURL
	usersShortURLs []model.Users_ShortURLs
}

func NewMockRepository() (*MockRepository, error) {
	users := []model.User{
		*model.NewUser(0, "User1"),
		*model.NewUser(0, "User2"),
		*model.NewUser(0, "User3"),
		*model.NewUser(0, "User4"),
	}
	urls := []model.ShortURL{
		*model.NewShortURL(0, "URL1", "Alias1"),
		*model.NewShortURL(0, "URL2", "Alias2"),
		*model.NewShortURL(0, "URL3", "Alias3"),
		*model.NewShortURL(0, "URL4", "Alias4"),
	}
	usersURL := []model.Users_ShortURLs{
		*model.NewUsers_ShortURLs(0, 0, 0),
		*model.NewUsers_ShortURLs(1, 1, 1),
		*model.NewUsers_ShortURLs(2, 2, 2),
		*model.NewUsers_ShortURLs(3, 3, 3),
	}

	return &MockRepository{
		userTable:      users,
		shortURLTable:  urls,
		usersShortURLs: usersURL,
	}, nil
}

func (repository MockRepository) GetAllShortURL() ([]model.ShortURL, error) { return nil, nil }

func (repository MockRepository) GetShortURLByID(id int) (*model.ShortURL, error) { return nil, nil }
func (repository MockRepository) GetShortURLByAlias(alias string) (*model.ShortURL, error) {
	return nil, nil
}
func (repository MockRepository) GetShortURLByURL(URL string) (*model.ShortURL, error) {
	return nil, nil
}

func (repository MockRepository) SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error) {
	return nil, nil
}

func (repository MockRepository) GetAllUsers() ([]model.User, error) { return nil, nil }

func (repository MockRepository) GetUserByID(ID int) (*model.User, error)        { return nil, nil }
func (repository MockRepository) GetUserByName(name string) (*model.User, error) { return nil, nil }

func (repository MockRepository) SaveUser(User model.User) (*model.User, error) { return nil, nil }

func (repository MockRepository) GetAllUsersShortURLS() ([]model.Users_ShortURLs, error) {
	return nil, nil
}

func (repository MockRepository) GetUsersShortURLSByID(ID int) (*model.Users_ShortURLs, error) {
	return nil, nil
}
func (repository MockRepository) GetUsersShortURLSByUserID(userID int) ([]model.Users_ShortURLs, error) {
	return nil, nil
}
func (repository MockRepository) GetUsersShortURLSByURLID(urlID int) (*model.Users_ShortURLs, error) {
	return nil, nil
}

func (repository MockRepository) SaveUserShortURL(shortURL model.Users_ShortURLs) (*model.Users_ShortURLs, error) {
	return nil, nil
}

func (repository MockRepository) CheckDBConnection() error { return nil }

func (repository MockRepository) Close() error { return nil }
