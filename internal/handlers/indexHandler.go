package handlers

import (
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/go-chi/chi/v5"
)

type indexHandler struct {
	config *config.Config
}

func NewIndexHandler(conf *config.Config) *indexHandler {
	return &indexHandler{
		config: conf,
	}
}

func (handler indexHandler) HandleGET(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	url, err := services.GetURL(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("location", string(url))
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func (handler indexHandler) HandlePOST(rw http.ResponseWriter, r *http.Request) {
	url, _ := io.ReadAll(r.Body)
	alias := services.CreateAlias(string(url))

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	_, err := rw.Write([]byte(handler.config.BaseURL + "/" + alias))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
