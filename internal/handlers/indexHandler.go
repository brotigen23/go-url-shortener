package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/brotigen23/go-url-shortener/internal/storage"
)

func IndexHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// -------------------------------------------------------------------------- POST
	case http.MethodPost:
		// Считывание
		url, _ := io.ReadAll(r.Body)
		alias := services.CreateAlias(string(url))

		// Заголовки и статус ответа
		rw.Header().Set("content-type", "text/plain")
		rw.WriteHeader(http.StatusCreated)

		// Запись ответа
		rw.Write([]byte(alias))
	// -------------------------------------------------------------------------- GET
	case http.MethodGet:
		alias := r.URL.Path[1:]
		url := services.GetURL(alias)
		if url == "" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("location", string(url))
		fmt.Println("Location: ", rw.Header().Get("location"))
		//http.Redirect(rw, r, string(url), http.StatusTemporaryRedirect)
		rw.Write([]byte(url))
	default:
		rw.WriteHeader(http.StatusNotAcceptable)
		rw.Write([]byte("Несуществующий метод"))

	}

	fmt.Println("Storage:")
	storage.Storage.Print()
}
