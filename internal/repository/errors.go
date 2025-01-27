package repository

import "errors"

var (
	ErrShortURLAlreadyExists = errors.New("short url already registered")
	ErrNoFound               = errors.New("short url not found")
)
