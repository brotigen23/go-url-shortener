package dto

type URLRequest struct {
	URL string `json:"url"`
}

func NewURLRequest() *URLRequest {
	return &URLRequest{}
}
