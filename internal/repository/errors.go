package repository

import "errors"

// Ошибки использующиеся репозиторием
var (
	ErrShortURLAlreadyExists = errors.New("short url already registered")
	ErrNoFound               = errors.New("short url not found")
)
