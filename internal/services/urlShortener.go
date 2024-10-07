package services

import (
	"github.com/brotigen23/go-url-shortener/internal/storage"
)

func CreateAlias(url string) string {
	return storage.Storage.Put([]byte(url))
}

func GetURL(alias string) string {
	return storage.Storage.FindByAlias([]byte(alias))
}
