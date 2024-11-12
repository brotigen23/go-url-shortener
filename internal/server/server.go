package server

// main.go

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handlers"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(conf *config.Config) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("logger init error: %v", err)
	}

	r := chi.NewRouter()

	aliases, _ := utils.LoadLocalAliases(conf.FileStoragePath)

	indexHandler := handlers.NewIndexHandler(conf, aliases)

	r.Get("/{id}", handlers.WithLogging(handlers.WithZip(indexHandler.HandleGET), logger.Sugar()))
	r.Get("/ping", handlers.WithLogging(handlers.WithZip(indexHandler.Ping), logger.Sugar()))

	r.Post("/", handlers.WithLogging(handlers.GzipMiddleware(indexHandler.HandlePOST), logger.Sugar()))

	r.Post("/api/shorten", handlers.WithLogging(handlers.GzipMiddleware(indexHandler.HandlePOSTAPI), logger.Sugar()))

	logger.Sugar().Infoln(
		"Server is running",
		"Server address", conf.ServerAddress,
		"Base URL", conf.BaseURL,
	)
	server := &http.Server{Addr: conf.ServerAddress, Handler: r}
	start := time.Now()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			return
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	duration := time.Since(start)

	logger.Sugar().Infoln(
		"server shutdown",
		"time running", duration,
	)
	err = utils.SaveLocalAliases(indexHandler.GetAliases(), conf.FileStoragePath)
	if err != nil {
		return err
	}
	return nil
}
