package dto

type APIShortenResponse struct {
	Result string `json:"result"`
}

func NewApiShortenResponse(result string) *APIShortenResponse {
	return &APIShortenResponse{
		Result: result,
	}
}
