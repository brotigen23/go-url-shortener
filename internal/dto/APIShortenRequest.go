package dto

type APIShortenRequest struct {
	URL string `json:"url"`
}

func NewAPIShortenRequest(URL string) *APIShortenRequest {
	return &APIShortenRequest{
		URL: URL,
	}
}
