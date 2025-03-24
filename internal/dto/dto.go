package dto

// DTO структура для запроса в виде JSON
type ShortenRequest struct {
	URL string `json:"url"`
}

// DTO структура для ответа в виде JSON
type ShortenResponse struct {
	Result string `json:"result"`
}

// DTO структура для запроса Batch в виде JSON
type BatchRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

// DTO структура для ответа Batch в виде JSON
type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// DTO структура для ответа в виде JSON со всеми сохраненными пользователем ссылками
type UserURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Number of urls and users
type Stats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}
