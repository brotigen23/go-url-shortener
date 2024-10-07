package handlers

import (
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/go-chi/chi/v5"
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
		rw.Write([]byte("http://localhost:8080/" + alias))
	// -------------------------------------------------------------------------- GET
	case http.MethodGet:
		alias := r.URL.Path[1:]
		url := services.GetURL(alias)
		if url == "" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("location", string(url))
		rw.WriteHeader(http.StatusTemporaryRedirect)
	default:
		rw.WriteHeader(http.StatusNotAcceptable)
		rw.Write([]byte("Несуществующий метод"))
	}
}

func IndexGET(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	url := services.GetURL(alias)
	if url == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("location", string(url))
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func IndexPOST(rw http.ResponseWriter, r *http.Request) {
	url, _ := io.ReadAll(r.Body)
	alias := services.CreateAlias(string(url))

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	rw.Write([]byte(config.ConfigENV.BastHostForAliases + "/" + alias))
}
