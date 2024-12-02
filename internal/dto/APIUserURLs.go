package dto

type APIUserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewAPIUserURLs(url, alias string) *APIUserURLs {
	return &APIUserURLs{
		ShortURL:    alias,
		OriginalURL: url,
	}
}
