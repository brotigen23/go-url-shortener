package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/repositories"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type mainHandler struct {
	config  *config.Config
	service *services.ServiceShortener
}

func NewMainHandler(conf *config.Config, aliases []model.ShortURL, logger *zap.Logger, repository repositories.Repository) (*mainHandler, error) {
	service, err := services.NewService(conf, 8, aliases, logger, repository)
	if err != nil {
		return nil, err
	}

	return &mainHandler{
		config:  conf,
		service: service,
	}, nil
}

// Store new ShortURL
func (handler *mainHandler) CreateShortURL(rw http.ResponseWriter, r *http.Request) {
	// ------------------------------- Read request data -------------------------------
	var URL string
	switch r.Header.Get("content-type") {
	case "text/plain":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		URL = string(body)
	case "application/json":
		request := &dto.APIShortenRequest{}
		var buffer bytes.Buffer
		_, err := buffer.ReadFrom(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buffer.Bytes(), &request); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		URL = request.URL
	}

	// ------------------------------- Save URL -------------------------------
	userName, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	alias, err := handler.service.SaveURL(userName.Value, URL)
	if err != nil {
		if err.Error() == "URL already exists" {
			rw.WriteHeader(http.StatusConflict)
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// ------------------------------- Write response -------------------------------
	rw.WriteHeader(http.StatusCreated)

	var response []byte

	switch r.Header.Get("content-type") {
	case "text/plain":
		rw.Header().Set("content-type", "text/plain")
		response = []byte(handler.config.BaseURL + "/" + alias)
	case "application/json":
		result := dto.NewAPIShortenResponse(handler.config.BaseURL + "/" + alias)
		response, err = json.Marshal(result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

// Store new ShortURLs
func (handler *mainHandler) CreateShortURLs(rw http.ResponseWriter, r *http.Request) {
	request := []dto.APIBatchRequest{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var URLs []string
	for _, url := range request {
		URLs = append(URLs, url.URL)
	}

	BatchResponse := []*dto.APIBatchResponse{}
	userName, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	shortURLs, err := handler.service.SaveURLs(userName.Value, URLs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	for key, value := range shortURLs {
		BatchResponse = append(BatchResponse, dto.NewAPIBatchResponse(key, value))
	}
	// Заголовки и статус ответа
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	response, err := json.Marshal(BatchResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

}

// Return URL by Alias
func (handler *mainHandler) GetShortURL(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	userName, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	URL, err := handler.service.GetURL(userName.Value, alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("location", URL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

// Return all user's URLs
func (handler *mainHandler) GetShortURLs(rw http.ResponseWriter, r *http.Request) {
	userName, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	URLs, err := handler.service.GetURLs(userName.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	if len(URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	batchResponse := []dto.APIUserURLs{}
	for key, value := range URLs {
		batchResponse = append(batchResponse, *dto.NewAPIUserURLs(key, value))
	}
	// Заголовки и статус ответа
	rw.Header().Set("content-type", "application/json")

	// Запись ответа
	response, err := json.Marshal(batchResponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	rw.WriteHeader(http.StatusOK)
}

func (handler *mainHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	if handler.service.CheckDBConnection() != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}