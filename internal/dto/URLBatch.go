package dto

type URL struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type BatchResponse struct {
	ID    string `json:"correlation_id"`
	Alias string `json:"short_url"`
}
