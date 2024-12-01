package model

type Users_ShortURLs struct {
	Id      int `json:"ID"`
	User_ID int `json:"User_ID"`
	URL_ID  int `json:"URL_ID"`
}

func NewUsers_ShortURLs(id int, User_ID int, URL_ID int) *Users_ShortURLs{
	return &Users_ShortURLs{
		Id: id,
		User_ID: User_ID,
		URL_ID: URL_ID,
	}
}