package model

type Alias struct {
	URL   string `json:"original_url"`
	Alias string `json:"short_url"`
}

func NewAlias(url string, alias string) *Alias {
	return &Alias{
		URL:   url,
		Alias: alias,
	}
}

func (a Alias) GetURL() string {
	return a.URL
}

func (a Alias) GetAlias() string {
	return a.Alias
}
