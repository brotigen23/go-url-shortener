package server

// main.go

import (
	"fmt"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func Run(conf *config.Config) error {
	r := chi.NewRouter()
	indexHandler := handlers.NewIndexHandler(conf)
	r.Get("/{id}", indexHandler.HandleGET)
	r.Post("/", indexHandler.HandlePOST)
	fmt.Printf("server running on %v\nbase url for alias is %v\n", conf.ServerAddress, conf.BaseURL)
	return http.ListenAndServe(conf.ServerAddress, r)
}
