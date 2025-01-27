package model

type ShortURL struct {
	ID        int
	URL       string
	ShortURL  string
	Username  string
	IsDeleted bool
}

func NewShortURLs(aliases []string) []ShortURL {
	return nil
}
