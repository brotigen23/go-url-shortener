package handler

import (
	"github.com/brotigen23/go-url-shortener/internal/middleware"
	"github.com/brotigen23/go-url-shortener/internal/repository/memory"
	"github.com/brotigen23/go-url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func ExampleHandler_CreateShortURL() {
	//------------------------------------------------------------
	// Create Repository
	//------------------------------------------------------------
	repo := memory.New(nil)

	//------------------------------------------------------------
	// Create service
	//------------------------------------------------------------
	serviceShortener, err := service.New(nil, logger, repo)
	if err != nil {
		return
	}

	handler, err := New("", serviceShortener)
	if err != nil {
		return
	}

	logger.Debugln("handler is initialized")

	//------------------------------------------------------------
	// Create mux with chi
	//------------------------------------------------------------
	r := chi.NewRouter()

	r.Use(middleware.Log(logger))
	r.Use(middleware.Auth("", logger))
	r.Use(middleware.Encoding)

	//------------------------------------------------------------
	// Use CreateShortURL
	//------------------------------------------------------------
	r.Post("/", handler.CreateShortURL)
}
