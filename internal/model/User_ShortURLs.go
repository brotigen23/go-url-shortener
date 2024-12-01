package model

type UsersShortURLs struct {
	ID     int `json:"ID"`
	UserID int `json:"User_ID"`
	UrlID  int `json:"URL_ID"`
}

func NewUsersShortURLs(id int, UserID int, URLID int) *UsersShortURLs {
	return &UsersShortURLs{
		ID:     id,
		UserID: UserID,
		UrlID:  URLID,
	}
}
