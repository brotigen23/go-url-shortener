package model


type Alias struct {
	url   string
	alias string
}

func NewAlias(url string, alias string) *Alias {
	return &Alias{
		url:   url,
		alias: alias,
	}
}

func (a Alias) GetURL() string {
	return a.url
}

func (a Alias) GetAlias() string {
	return a.alias
}
