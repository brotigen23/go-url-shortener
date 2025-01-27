package model

type ShortURL struct {
	ID        int
	URL       string
	Alias     string
	Username  string
	IsDeleted bool
}

func NewShortURLs(aliases []string) []ShortURL {
	return nil
}
