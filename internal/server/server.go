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
	"github.com/brotigen23/go-url-shortener/internal/repositories"
	"github.com/brotigen23/go-url-shortener/internal/services"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(conf *config.Config) error {
	logger, err := zap.NewDevelopment()
	var repository repositories.Repository
	if conf.DatabaseDSN != "" {
		repository, err = repositories.NewPostgresRepository("postgres", conf.DatabaseDSN, logger)
		if err != nil {
			return err
		}
	} else {
		repository = repositories.NewInMemoryRepo(nil, nil, nil)
	}

	if err != nil {
		return fmt.Errorf("logger init error: %v", err)
	}

	r := chi.NewRouter()

	authService, err := services.NewServiceAuth(conf, repository)
	if err != nil {
		return err
	}

	aliases, _ := utils.LoadLocalAliases(conf.FileStoragePath)

	indexHandler, err := handlers.NewIndexHandler(conf, aliases, logger, repository)
	if err != nil {
		return err
	}

	mainHandler, err := handlers.NewMainHandler(conf, aliases, logger, repository)
	if err != nil {
		return err
	}

	//r.Get("/{id}", handlers.WithLogging(handlers.WithAuth(handlers.WithZip(indexHandler.HandleGET), conf, authService), logger.Sugar()))
	r.Get("/{id}", handlers.WithLogging(handlers.WithAuth(handlers.WithZip(mainHandler.GetShortURL), conf, authService), logger.Sugar()))
	r.Get("/ping", handlers.WithLogging(handlers.WithAuth(handlers.WithZip(indexHandler.Ping), conf, authService), logger.Sugar()))
	//r.Get("/api/user/urls", handlers.WithLogging(handlers.WithAuth(handlers.WithZip(indexHandler.GetUsersURL), conf, authService), logger.Sugar()))
	r.Get("/api/user/urls", handlers.WithLogging(handlers.WithAuth(handlers.WithZip(mainHandler.GetShortURLs), conf, authService), logger.Sugar()))

	//r.Post("/", handlers.WithLogging(handlers.WithAuth(handlers.GzipMiddleware(indexHandler.HandlePOST), conf, authService), logger.Sugar()))
	r.Post("/", handlers.WithLogging(handlers.WithAuth(handlers.GzipMiddleware(mainHandler.CreateShortURL), conf, authService), logger.Sugar()))
	//r.Post("/api/shorten", handlers.WithLogging(handlers.WithAuth(handlers.GzipMiddleware(indexHandler.HandlePOSTAPI), conf, authService), logger.Sugar()))
	r.Post("/api/shorten", handlers.WithLogging(handlers.WithAuth(handlers.GzipMiddleware(mainHandler.CreateShortURL), conf, authService), logger.Sugar()))
	//r.Post("/api/shorten/batch", handlers.WithLogging(handlers.WithAuth(handlers.GzipMiddleware(indexHandler.Batch), conf, authService), logger.Sugar()))
	r.Post("/api/shorten/batch", handlers.WithLogging(handlers.WithAuth(handlers.GzipMiddleware(mainHandler.CreateShortURLs), conf, authService), logger.Sugar()))

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
	if conf.DatabaseDSN == "" {
		err = utils.SaveStorage(indexHandler.GetAliases(), conf.FileStoragePath)
		if err != nil {
			return err
		}
	}
	return nil
}
