package server

// main.go

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handler"
	"github.com/brotigen23/go-url-shortener/internal/middleware"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/repository/memory"
	"github.com/brotigen23/go-url-shortener/internal/repository/postgres"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Run(config *config.Config, logger *zap.SugaredLogger) error {

	//------------------------------------------------------------
	// REPOSITORY
	//------------------------------------------------------------
	var repo repository.Repository
	const driver = "postgres"
	switch config.DatabaseDSN {
	case "":
		repo = memory.New(nil)
	default:
		db, err := sql.Open(driver, config.DatabaseDSN)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			return err
		}
		repo = postgres.New(db, logger)
	}
	logger.Debugln("repo is initialized")

	//------------------------------------------------------------
	// SERVICES AND HANDLER
	//------------------------------------------------------------
	serviceShortener, err := service.New(config, logger, repo)
	if err != nil {
		return err
	}
	logger.Debugln("service is initialized")

	handler, err := handler.New(config.BaseURL, serviceShortener)
	if err != nil {
		return err
	}

	logger.Debugln("handler is initialized")

	//------------------------------------------------------------
	// ROUTER
	//------------------------------------------------------------
	r := chi.NewRouter()

	r.Use(middleware.Log(logger))
	r.Use(middleware.Auth(config.JWTSecretKey, logger)) // TODO: secret key from config
	r.Use(middleware.Encoding)

	r.Get("/{id}", handler.GetShortURL)
	r.Get("/ping", handler.Ping)
	r.Get("/api/user/urls", handler.GetShortURLs)
	r.Delete("/api/user/urls", handler.Detele)
	r.Post("/", handler.CreateShortURL)
	r.Post("/api/shorten", handler.CreateShortURL)
	r.Post("/api/shorten/batch", handler.CreateShortURLs)

	logger.Debugln("router is initialized")

	//------------------------------------------------------------
	// SERVER
	//------------------------------------------------------------
	server := &http.Server{Addr: config.ServerAddress, Handler: r}
	start := time.Now()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			return
		}
	}()
	logger.Infoln("Server is running")

	//------------------------------------------------------------
	// GRACEFULL SHUTDOWN
	//------------------------------------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	duration := time.Since(start)

	logger.Infoln(
		"server shutdown",
		"time running", duration,
	)
	shortURLs, err := repo.GetAll()
	if err != nil {
		return err
	}
	err = utils.SaveStorage(shortURLs, config.FileStoragePath)
	if err != nil {
		return err
	}
	return nil
}
