package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/go-chi/chi/v5"
)

type IndexHandler struct {
	config  *config.Config
	service *services.ServiceShortener
}

func NewIndexHandler(conf *config.Config, aliases []model.Alias) *IndexHandler {
	return &IndexHandler{
		config:  conf,
		service: services.NewService(8, aliases),
	}
}

func (handler IndexHandler) HandleGET(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	URL, err := handler.service.GetURLByAlias(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("location", URL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func (handler IndexHandler) HandlePOST(rw http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println("BODY: ", string(body))
	alias, err := handler.service.Save(string(body))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	_, err = rw.Write([]byte(handler.config.BaseURL + "/" + alias))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}

func (handler IndexHandler) HandlePOSTAPI(rw http.ResponseWriter, r *http.Request) {
	req := dto.NewURLRequest()
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	alias, err := handler.service.Save(req.URL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	resp := dto.NewAliasResponse()
	resp.Result = handler.config.BaseURL + "/" + alias

	// Запись ответа
	response, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}

func (handler IndexHandler) GetAliases() []model.Alias {
	return *handler.service.GetAll()
}
