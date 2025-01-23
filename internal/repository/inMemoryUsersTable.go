package repository

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

func (r *inMemoryRepo) GetAllUsers() ([]model.User, error) { return nil, nil }

func (r *inMemoryRepo) GetUserByID(ID int) (*model.User, error) {
	for _, user := range r.Users {
		if user.ID == ID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
func (r *inMemoryRepo) GetUserByName(name string) (*model.User, error) {
	for _, user := range r.Users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *inMemoryRepo) SaveUser(User model.User) (*model.User, error) {
	user := model.NewUser(len(r.Users), User.Name)
	r.Users = append(r.Users, *user)
	return user, nil
}
