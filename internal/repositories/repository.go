package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

type Short_URLs interface {
	GetAllShortURL() ([]model.ShortURL, error)

	GetShortURLByID(id int) (*model.ShortURL, error)
	GetShortURLByAlias(alias string) (*model.ShortURL, error)
	GetShortURLByURL(URL string) (*model.ShortURL, error)

	SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error)
}

type Users interface {
	GetAllUsers() ([]model.User, error)

	GetUserByID(ID int) (*model.User, error)
	GetUserByName(name string) (*model.User, error)

	SaveUser(User model.User) (*model.User, error)
}

type Users_ShortURLs interface {
	GetAllUsersShortURLS() ([]model.Users_ShortURLs, error)

	GetUsersShortURLSByID(ID int) (*model.Users_ShortURLs, error)
	GetUsersShortURLSByUserID(userID int) ([]model.Users_ShortURLs, error)
	GetUsersShortURLSByURLID(urlID int) (*model.Users_ShortURLs, error)

	SaveUserShortURL(shortURL model.Users_ShortURLs) (*model.Users_ShortURLs, error)
}

type Repository interface {
	Short_URLs
	Users
	Users_ShortURLs

	CheckDBConnection() error

	Close() error
}

type Repo interface {
	GetByAlias(alias string) (*model.ShortURL, error)
	GetByURL(url string) (*model.ShortURL, error)
	GetAll() *[]model.ShortURL
	Save(model model.ShortURL) error
	Migrate(model []model.ShortURL)
	CheckDBConnection() error
	SaveUser(id string) error
	GetUserByID(userID string) error
	SaveUserURL(userID string, alias string) error
	GetUserURL(userID string) ([]model.ShortURL, error)

	Close()
}
