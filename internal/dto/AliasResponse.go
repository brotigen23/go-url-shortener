package dto

type AliasResponse struct {
	Result string `json:"result"`
}

func NewAliasResponse() *AliasResponse {
	return &AliasResponse{}
}
