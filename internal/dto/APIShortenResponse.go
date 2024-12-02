package dto

type APIShortenResponse struct {
	Result string `json:"result"`
}

func NewAPIShortenResponse(result string) *APIShortenResponse {
	return &APIShortenResponse{
		Result: result,
	}
}
