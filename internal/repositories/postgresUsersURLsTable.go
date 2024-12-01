package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

//---------------------- Users_ShortURLs table ----------------------

func (repo PostgresRepository) GetAllUsersShortURLS() ([]model.Users_ShortURLs, error) {
	return nil, nil
}

func (repo PostgresRepository) GetUsersShortURLSByID(ID int) (*model.Users_ShortURLs, error) {
	return nil, nil
}

// Return all user's shortURL by UserID
func (repo PostgresRepository) GetUsersShortURLSByUserID(userID int) ([]model.Users_ShortURLs, error) {
	//query := "SELECT * FROM Users_URLs WHERE ID IN (( SELECT URL_ID FROM Users_URLs WHERE User_ID = $1))"
	query := "SELECT * FROM Users_URLs WHERE User_ID = $1"
	q, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var id int
	var UserID int
	var URLID int
	var ret []model.Users_ShortURLs
	for q.Next() {
		err = q.Scan(&id, &UserID, &URLID)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *model.NewUsers_ShortURLs(id, UserID, URLID))
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ret, nil
}

// Return all users by ShortURL
func (repo PostgresRepository) GetUsersShortURLSByURLID(urlID int) (*model.Users_ShortURLs, error) {
	return nil, nil
}

func (repo PostgresRepository) SaveUserShortURL(Users_ShortURLs model.Users_ShortURLs) (*model.Users_ShortURLs, error) {
	query := "INSERT INTO Users_URLs(User_ID, URL_ID) VALUES($1, $2) RETURNING ID"
	var (
		id int
	)
	err := repo.db.QueryRow(query, Users_ShortURLs.User_ID, Users_ShortURLs.URL_ID).Scan(&id)
	if err != nil {
		return nil, err
	}
	return model.NewUsers_ShortURLs(id, Users_ShortURLs.User_ID, Users_ShortURLs.URL_ID), nil
}
