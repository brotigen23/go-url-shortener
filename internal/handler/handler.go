package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/database"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	baseURL string

	service *service.Service
}

func New(baseURL string, service *service.Service) (*handler, error) {
	return &handler{
		service: service,
		baseURL: baseURL,
	}, nil
}

// Store new ShortURL
func (h *handler) CreateShortURL(rw http.ResponseWriter, r *http.Request) {

	// --------------------------------------------------------------
	// READ REQUEST DATA
	// --------------------------------------------------------------
	var URL string
	switch r.Header.Get("content-type") {
	case "text/plain; charset=utf-8", "application/x-gzip", "text/plain":
		rw.Header().Set("content-type", "text/plain")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		URL = string(body)
	case "application/json":
		rw.Header().Set("content-type", "application/json")
		request := &dto.ShortenRequest{}
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

	if URL == "" {
		http.Error(rw, "url is empty", http.StatusBadRequest)
		return
	}
	// ------------------------------- Save URL -------------------------------
	userName, err := r.Cookie("username")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	alias, err := h.service.CreateShortURL(userName.Value, URL)
	if err != nil {
		if err == service.ErrShortURLAlreadyExists {
			http.Error(rw, "short url already saved", http.StatusConflict)
			return
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// ------------------------------- Write response -------------------------------
	var response []byte

	switch r.Header.Get("content-type") {
	case "text/plain; charset=utf-8", "text/plain", "application/x-gzip":
		response = []byte(h.baseURL + "/" + alias)
	case "application/json":
		result := dto.ShortenResponse{Result: h.baseURL + "/" + alias}
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
func (h *handler) CreateShortURLs(rw http.ResponseWriter, r *http.Request) {
	request := []dto.BatchRequest{}
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

	BatchResponse := []*dto.BatchResponse{}
	userName, err := r.Cookie("username")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	shortURLs, err := h.service.CreateShortURLs(userName.Value, URLs)
	if err != nil {
		if err == service.ErrShortURLAlreadyExists {
			http.Error(rw, "some short url already register", http.StatusConflict)
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// TODO: перенести создание Batch Response в коструктор с параметром map[string]string
	for i := range request {
		BatchResponse = append(BatchResponse, &dto.BatchResponse{ID: request[i].ID, ShortURL: h.baseURL + "/" + shortURLs[request[i].URL]})
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
func (h *handler) GetShortURL(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	URL, err := h.service.GetShortURL(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	isDeleted, err := h.service.IsShortURLDeleted(alias)
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

func (h *handler) GetShortURLs(rw http.ResponseWriter, r *http.Request) {
	userName, err := r.Cookie("username")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	URLs, err := h.service.GetShortURLs(userName.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	if len(URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	batchResponse := []dto.UserURLs{}
	for key, value := range URLs {
		batchResponse = append(batchResponse, dto.UserURLs{OriginalURL: key, ShortURL: h.baseURL + "/" + value})
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

func (h *handler) Detele(rw http.ResponseWriter, r *http.Request) {
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
	URLs, err := h.service.GetShortURLs(userName.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if len(URLs) == 0 {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Println(request)
	err = h.service.DeleteShortURLs(userName.Value, request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}

func (h *handler) Ping(rw http.ResponseWriter, r *http.Request) {
	if err := database.CheckPostgresConnection(h.service.GetDSN()); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
