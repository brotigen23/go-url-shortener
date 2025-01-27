package dto

// Create short URL json
type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

// Batch
type BatchRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// Get User's URLs
type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
