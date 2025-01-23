package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type mainHandler struct {
	service *service.ServiceShortener
}

func NewMainHandler(logger *zap.Logger, service *service.ServiceShortener) (*mainHandler, error) {
	return &mainHandler{
		service: service,
	}, nil
}

// Store new ShortURL
func (h *mainHandler) CreateShortURL(rw http.ResponseWriter, r *http.Request) {
    switch r.Header.Get("content-type") {
	case "text/plain; charset=utf-8", "text/plain", "application/x-gzip":
		rw.Header().Set("content-type", "text/plain")
	case "application/json":
		rw.Header().Set("content-type", "application/json")
	}

	// ------------------------------- Read request data -------------------------------
	var URL string
	switch r.Header.Get("content-type") {
	case "text/plain; charset=utf-8", "text/plain", "application/x-gzip":
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
	alias, err := h.service.SaveURL(userName.Value, URL)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "short_urls_url_key"` {
			rw.WriteHeader(http.StatusConflict)
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// ------------------------------- Write response -------------------------------
	var response []byte

	switch r.Header.Get("content-type") {
	case "text/plain; charset=utf-8", "text/plain", "application/x-gzip":
		response = []byte(h.service.GetBaseURL() + "/" + alias)
	case "application/json":
		result := dto.NewAPIShortenResponse(h.service.GetBaseURL() + "/" + alias)
		response, err = json.Marshal(result)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

// Store new ShortURLs
func (h *mainHandler) CreateShortURLs(rw http.ResponseWriter, r *http.Request) {
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
	shortURLs, err := h.service.SaveURLs(userName.Value, URLs)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "short_urls_url_key"` {
			rw.WriteHeader(http.StatusConflict)
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	for i := range request {
		BatchResponse = append(BatchResponse, dto.NewAPIBatchResponse(request[i].ID, h.service.GetBaseURL()+"/"+shortURLs[request[i].URL]))
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
func (h *mainHandler) GetShortURL(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	URL, err := h.service.GetURL(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	isDeleted, err := h.service.IsDeleted(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if isDeleted {
		rw.WriteHeader(http.StatusGone)
		return
	}
	rw.Header().Set("location", URL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *mainHandler) GetShortURLs(rw http.ResponseWriter, r *http.Request) {
	userName, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	URLs, err := h.service.GetURLs(userName.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	if len(URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	batchResponse := []dto.APIUserURLs{}
	for key, value := range URLs {
		batchResponse = append(batchResponse, *dto.NewAPIUserURLs(key, h.service.GetBaseURL()+"/"+value))
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

func (h *mainHandler) Detele(rw http.ResponseWriter, r *http.Request) {
	request := []string{}
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
	userName, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	URLs, err := h.service.GetURLs(userName.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if len(URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Println(request)
	err = h.service.DeleteURLs(userName.Value, request)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}

func (h *mainHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	if h.service.CheckDBConnection() != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
