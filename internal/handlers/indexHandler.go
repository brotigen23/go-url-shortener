package handlers

import (
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/go-chi/chi/v5"
)

type indexHandler struct {
	config  *config.Config
	service *services.ServiceShortener
}

func NewIndexHandler(conf *config.Config) *indexHandler {
	return &indexHandler{
		config:  conf,
		service: services.NewService(8),
	}
}

func (handler indexHandler) HandleGET(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	model, err := handler.service.GetURLByAlias(alias) //handler.repo.GetByAlias(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("location", string(model.GetURL()))
	rw.WriteHeader(http.StatusTemporaryRedirect)
}


func (handler indexHandler) HandlePOST(rw http.ResponseWriter, r *http.Request) {
	url, _ := io.ReadAll(r.Body)

	model, err := handler.service.Save(string(url))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	_, err = rw.Write([]byte(handler.config.BaseURL + "/" + model.GetAlias()))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
