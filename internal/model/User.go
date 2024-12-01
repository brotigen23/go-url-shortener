package model

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewUser(id int, name string) *User {
	return &User{
		Id:   id,
		Name: name,
	}
}
