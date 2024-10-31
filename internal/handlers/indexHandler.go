package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
)

/*
Добавьте в код сервера новый эндпоинт POST /api/shorten, который будет принимать в теле запроса JSON-объект
{"url":"<some_url>"} и возвращать в ответ объект {"result":"<short_url>"}.
Запрос может иметь такой вид:

POST http://localhost:8080/api/shorten HTTP/1.1
Host: localhost:8080
Content-Type: application/json
{
  "url": "https://practicum.yandex.ru"
}

Ответ может быть таким:

HTTP/1.1 201 OK
Content-Type: application/json
Content-Length: 30
{
 "result": "http://localhost:8080/EwHXdJfB"
}
*/

type req struct {
	URL string `json:"url"`
}

type resp struct {
	Result string `json:"result"`
}

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
func NewMockIndexHandler(conf *config.Config) *indexHandler {
	return &indexHandler{
		config:  conf,
		service: services.NewMockService(8),
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

	var url []byte
	if r.Header.Get("Content-Encoding") == "gzip" {
		url, _ = utils.Unzip(r.Body)
	} else {
		url, _ = io.ReadAll(r.Body)
	}

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

func (handler indexHandler) HandlePOSTAPI(rw http.ResponseWriter, r *http.Request) {
	req := new(req)
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

	model, err := handler.service.Save(req.URL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	resp := new(resp)
	resp.Result = handler.config.BaseURL + "/" + model.GetAlias()

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
