package repository

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

//---------------------- Users table ----------------------

func (r PostgresRepository) GetAllUsers() ([]model.User, error) { return nil, nil }

func (r PostgresRepository) GetUserByID(ID int) (*model.User, error) { return nil, nil }
func (r PostgresRepository) GetUserByName(name string) (*model.User, error) {
	fmt.Println(name)
	query := r.db.QueryRow(`SELECT * FROM Users WHERE Name = $1`, name)
	var ID int
	var Name string
	err := query.Scan(&ID, &Name)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: ID, Name: Name}, nil
}

func (r PostgresRepository) SaveUser(User model.User) (*model.User, error) {
	query := "INSERT INTO Users(Name) VALUES($1) RETURNING ID"
	var (
		id int
	)
	err := r.db.QueryRow(query, User.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	return model.NewUser(id, User.Name), nil
}
