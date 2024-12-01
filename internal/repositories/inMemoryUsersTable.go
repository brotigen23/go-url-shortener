package repositories

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/model"
)

func (repo *inMemoryRepo) GetAllUsers() ([]model.User, error) { return nil, nil }

func (repo *inMemoryRepo) GetUserByID(ID int) (*model.User, error) {
	for _, user := range repo.Users {
		fmt.Println(user)
		if user.Id == ID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
func (repo *inMemoryRepo) GetUserByName(name string) (*model.User, error) {
	for _, user := range repo.Users {
		fmt.Println(user)
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (repo *inMemoryRepo) SaveUser(User model.User) (*model.User, error) {
	user := model.NewUser(len(repo.Users), User.Name)
	fmt.Println(User)
	repo.Users = append(repo.Users, *user)
	return user, nil
}
