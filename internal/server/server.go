package server

// main.go

import (
	"fmt"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handlers"
	"github.com/brotigen23/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Run(conf *config.Config, stor *storage.Storage) error {
	r := chi.NewRouter()
	indexHandler := handlers.NewIndexHandler(conf, stor)
	r.Get("/{id}", indexHandler.HandleGET)
	r.Post("/", indexHandler.HandlePOST)
	fmt.Printf("server is running on %v\nbase url for alias is %v\n", conf.ServerAddress, conf.BaseURL)
	return http.ListenAndServe(conf.ServerAddress, r)
}
