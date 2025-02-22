package service

import "errors"

// Ошибки для использования сервисом
var (
	ErrShortURLAlreadyExists = errors.New("short url already exists")
	ErrShortURLNotFound      = errors.New("short url not found")
)
