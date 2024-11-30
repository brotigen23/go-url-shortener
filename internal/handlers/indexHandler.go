package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/dto"
	"github.com/brotigen23/go-url-shortener/internal/model"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type IndexHandler struct {
	config  *config.Config
	service *services.ServiceShortener
}

func NewIndexHandler(conf *config.Config, aliases []model.Alias, logger *zap.Logger) (*IndexHandler, error) {
	service, err := services.NewService(conf, 8, aliases, logger)
	if err != nil {
		return nil, err
	}
	return &IndexHandler{
		config:  conf,
		service: service,
	}, nil
}

func (handler IndexHandler) HandleGET(rw http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")
	user, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	URL, err := handler.service.GetURLByAlias(alias)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.Header().Set("location", URL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func (handler IndexHandler) HandlePOST(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "text/plain")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	alias, err := handler.service.Save(string(body))
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "aliases_url_key"` {
			rw.WriteHeader(http.StatusConflict)
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	user, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	handler.service.SaveUsersURL(user.Value, alias)
	// Заголовки и статус ответа
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	_, err = rw.Write([]byte(handler.config.BaseURL + "/" + alias))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}

func (handler IndexHandler) HandlePOSTAPI(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
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
		if err.Error() == `pq: duplicate key value violates unique constraint "aliases_url_key"` {
			rw.WriteHeader(http.StatusConflict)
		} else {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	user, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	handler.service.SaveUsersURL(user.Value, alias)

	// Заголовки и статус ответа
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

func (handler IndexHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	if handler.service.CheckDBConnection() != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)

}

func (handler IndexHandler) Batch(rw http.ResponseWriter, r *http.Request) {
	req := []dto.URL{}
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
	res := []*dto.BatchResponse{}
	user, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	for _, value := range req {
		alias, err := handler.service.Save(value.URL)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.service.SaveUsersURL(user.Value, alias)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		res = append(res, &dto.BatchResponse{ID: value.ID, Alias: handler.config.BaseURL + "/" + alias})
	}

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	response, err := json.Marshal(res)
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

func (handler IndexHandler) GetUsersURL(rw http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("userID")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	res, err := handler.service.GetUserURL(user.Value)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	// Заголовки и статус ответа
	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	// Запись ответа
	response, err := json.Marshal(res)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = rw.Write(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
