package model

type ShortURL struct {
	ID    int    `json:"-"`
	URL   string `json:"original_url"`
	Alias string `json:"short_url"`
	IsDeleted bool `json:"Is_Deleted"`
}

func NewShortURL(id int, url string, alias string) *ShortURL {
	return &ShortURL{
		ID:    id,
		URL:   url,
		Alias: alias,
	}
}
