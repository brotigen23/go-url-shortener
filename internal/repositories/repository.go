package repositories

import "github.com/brotigen23/go-url-shortener/internal/model"

type ShortURLs interface {
	GetAllShortURL() ([]model.ShortURL, error)

	GetShortURLByID(id int) (*model.ShortURL, error)
	GetShortURLByAlias(alias string) (*model.ShortURL, error)
	GetShortURLByURL(URL string) (*model.ShortURL, error)

	DeleteShortURLByAlias(Alias string) error
	DeleteShortURLByAliases(Aliases []string) error

	SaveShortURL(ShortURL model.ShortURL) (*model.ShortURL, error)
}

type Users interface {
	GetAllUsers() ([]model.User, error)

	GetUserByID(ID int) (*model.User, error)
	GetUserByName(name string) (*model.User, error)

	SaveUser(User model.User) (*model.User, error)
}

type UsersShortURLs interface {
	GetAllUsersShortURLS() ([]model.UsersShortURLs, error)

	GetUsersShortURLSByID(ID int) (*model.UsersShortURLs, error)
	GetUsersShortURLSByUserID(userID int) ([]model.UsersShortURLs, error)
	GetUsersShortURLSByURLID(urlID int) (*model.UsersShortURLs, error)

	SaveUserShortURL(shortURL model.UsersShortURLs) (*model.UsersShortURLs, error)
}

type Repository interface {
	ShortURLs
	Users
	UsersShortURLs

	CheckDBConnection() error

	Close() error
}
