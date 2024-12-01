package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

//---------------------- Users_ShortURLs table ----------------------

func (repo PostgresRepository) GetAllUsersShortURLS() ([]model.UsersShortURLs, error) {
	return nil, nil
}

func (repo PostgresRepository) GetUsersShortURLSByID(ID int) (*model.UsersShortURLs, error) {
	return nil, nil
}

// Return all user's shortURL by UserID
func (repo PostgresRepository) GetUsersShortURLSByUserID(userID int) ([]model.UsersShortURLs, error) {
	//query := "SELECT * FROM Users_URLs WHERE ID IN (( SELECT URL_ID FROM Users_URLs WHERE User_ID = $1))"
	query := "SELECT * FROM Users_URLs WHERE User_ID = $1"
	q, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var id int
	var UserID int
	var URLID int
	var ret []model.UsersShortURLs
	for q.Next() {
		err = q.Scan(&id, &UserID, &URLID)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *model.NewUsersShortURLs(id, UserID, URLID))
	}
	if q.Err() != nil {
		return nil, err
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

// Return all users by ShortURL
func (repo PostgresRepository) GetUsersShortURLSByURLID(urlID int) (*model.UsersShortURLs, error) {
	return nil, nil
}

func (repo PostgresRepository) SaveUserShortURL(UsersShortURLs model.UsersShortURLs) (*model.UsersShortURLs, error) {
	query := "INSERT INTO Users_URLs(User_ID, URL_ID) VALUES($1, $2) RETURNING ID"
	var (
		id int
	)
	err := repo.db.QueryRow(query, UsersShortURLs.UserID, UsersShortURLs.URLID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return model.NewUsersShortURLs(id, UsersShortURLs.UserID, UsersShortURLs.URLID), nil
}
