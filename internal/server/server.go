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

	r.Post("/", handlers.WithLogging(handlers.Withgzip(indexHandler.HandlePOST), logger.Sugar()))

	r.Post("/api/shorten", handlers.WithLogging(handlers.Withgzip(indexHandler.HandlePOSTAPI), logger.Sugar()))

	logger.Sugar().Infoln(
		"Server is running",
		"Server address", conf.ServerAddress,
		"Base URL", conf.BaseURL,
	)
	return http.ListenAndServe(conf.ServerAddress, r)
}
