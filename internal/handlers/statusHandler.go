package handlers

import (
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/storage"
)

func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	// намеренно добавлена ошибка в JSON
	rw.Write([]byte(storage.Storage.String()))
}
