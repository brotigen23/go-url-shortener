package model

type ShortURL struct {
	ID    int    `json:"-"`
	URL   string `json:"original_url"`
	Alias string `json:"short_url"`
}

func NewShortURL(id int, url string, alias string) *ShortURL {
	return &ShortURL{
		ID:    id,
		URL:   url,
		Alias: alias,
	}
}

func (a ShortURL) GetURL() string {
	return a.URL
}

func (a ShortURL) GetAlias() string {
	return a.Alias
}
