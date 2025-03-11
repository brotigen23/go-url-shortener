package server

// main.go

import (
	"context"
	"database/sql"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/database/migration"
	"github.com/brotigen23/go-url-shortener/internal/handler"
	"github.com/brotigen23/go-url-shortener/internal/middleware"
	"github.com/brotigen23/go-url-shortener/internal/repository"
	"github.com/brotigen23/go-url-shortener/internal/repository/memory"
	"github.com/brotigen23/go-url-shortener/internal/repository/postgres"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/brotigen23/go-url-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Производит конфигурирование, запуск и завершение работы сервера
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
		err = migration.MigratePostgresUp(db)
		if err != nil {
			return err
		}
		repo = postgres.New(db, logger)
	}
	//------------------------------------------------------------
	// SERVICES AND HANDLER
	//------------------------------------------------------------
	serviceShortener := service.New(config, logger, repo)

	handler := handler.New(config.BaseURL, serviceShortener)

	middleware := middleware.New(config.JWTSecretKey, logger)

	//------------------------------------------------------------
	// ROUTER
	//------------------------------------------------------------
	r := chi.NewRouter()

	r.Use(middleware.Log)
	r.Use(middleware.Auth) // TODO: secret key from config
	r.Use(middleware.Encoding)

	r.Get("/{id}", handler.RedirectByShortURL)
	r.Get("/ping", handler.Ping)
	r.Get("/api/user/urls", handler.GetShortURLs)

	r.Get("/stats", handler.Stats)

	r.Delete("/api/user/urls", handler.Detele)
	r.Post("/", handler.CreateShortURL)
	r.Post("/api/shorten", handler.CreateShortURL)
	r.Post("/api/shorten/batch", handler.CreateShortURLs)
	//------------------------------------------------------------
	// pprof
	//------------------------------------------------------------
	r.Mount("/debug", http.DefaultServeMux)

	logger.Debugln("router is initialized")

	//------------------------------------------------------------
	// SERVER
	//------------------------------------------------------------
	server := &http.Server{Addr: config.ServerAddress, Handler: r}
	start := time.Now()
	go func() {
		var err error
		if config.EnableHTTPS {
			c, k, er := utils.CreateCert()
			if er != nil {
				logger.Errorln(err)
				return
			}
			err = server.ListenAndServeTLS(c.String(), k.String())
		} else {
			err = server.ListenAndServe()
		}
		if err != nil {
			logger.Errorln(err)
			return
		}
	}()
	logger.Infoln("Server is running")

	//------------------------------------------------------------
	// GRACEFULL SHUTDOWN
	//------------------------------------------------------------
	stop := make(chan os.Signal, 4)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGQUIT)

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
	if config.DatabaseDSN == "" {
		shortURLs, err := repo.GetAll()
		if err != nil {
			return err
		}
		err = utils.SaveStorage(shortURLs, config.FileStoragePath)
		if err != nil {
			return err
		}
		logger.Infoln("short urls saved to", config.FileStoragePath)
	}
	return nil
}
