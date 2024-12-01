package dto

type APIBatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func NewAPIBatchResponse(id, shortURL string) *APIBatchResponse {
	return &APIBatchResponse{
		ID:       id,
		ShortURL: shortURL,
	}
}
