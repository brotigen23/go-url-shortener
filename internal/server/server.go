package server

// main.go

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handler"
	"github.com/brotigen23/go-url-shortener/internal/middleware"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run(config *config.Config) error {
	// Init
	//Logger
	logConf := zap.NewDevelopmentConfig()
	logConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := logConf.Build()
	if err != nil {
		return err
	}

	logger.Info("Now logs should be colored")
	var repo repository.Repository
	// Repo
	switch config.DatabaseDSN {
	case "":
		repo = repository.NewInMemoryRepo(nil, nil, nil)
	default:
		repo, err = repository.NewPostgresRepository("postgres", config.DatabaseDSN, logger)
		if err != nil {
			return err
		}
	}
	// Services

	serviceShortener, err := service.NewServiceShortener(config, 8, logger.Sugar(), repo)
	if err != nil {
		return err
	}
	authService, err := service.NewServiceAuth(config, repo)
	if err != nil {
		return err
	}

	// Handler
	mainHandler, err := handler.NewMainHandler(logger, serviceShortener)
	if err != nil {
		return err
	}

	// Router
	r := chi.NewRouter()

	r.Use(middleware.Log(logger.Sugar()))
	r.Use(middleware.Auth(config, authService, logger.Sugar()))
	r.Use(middleware.Encoding)

	r.Get("/{id}", mainHandler.GetShortURL)
	r.Get("/ping", mainHandler.Ping)
	r.Get("/api/user/urls", mainHandler.GetShortURLs)
	r.Delete("/api/user/urls", mainHandler.Detele)
	r.Post("/", mainHandler.CreateShortURL)
	r.Post("/api/shorten", mainHandler.CreateShortURL)
	r.Post("/api/shorten/batch", mainHandler.CreateShortURLs)

	// Logs
	logger.Sugar().Infoln(
		"Server is running",
		"Server address", config.ServerAddress,
		"Base URL", config.BaseURL,
	)

	// Server
	server := &http.Server{Addr: config.ServerAddress, Handler: r}
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
	shortURLs, err := repo.GetAllShortURL()
	if err != nil {
		return err
	}
	err = utils.SaveStorage(shortURLs, config.FileStoragePath)
	if err != nil {
		return err
	}
	return nil
}
