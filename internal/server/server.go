package server

// main.go

import (
	"fmt"
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(conf *config.Config) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("logger init error: %v", err)
	}

	r := chi.NewRouter()
	indexHandler := handlers.NewIndexHandler(conf)
	r.Get("/{id}", handlers.WithLogging(indexHandler.HandleGET, logger.Sugar()))
	r.Post("/", handlers.WithLogging(indexHandler.HandlePOST, logger.Sugar()))
	r.Post("/api/shorten", handlers.WithLogging(indexHandler.HandlePOSTAPI, logger.Sugar()))
	fmt.Printf("server is running on %v\nbase url for alias is %v\n", conf.ServerAddress, conf.BaseURL)
	return http.ListenAndServe(conf.ServerAddress, r)
}
