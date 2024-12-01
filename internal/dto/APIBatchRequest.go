package dto

type APIBatchRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}