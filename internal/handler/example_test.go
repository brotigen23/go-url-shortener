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
	serviceShortener := service.New(nil, logger, repo)

	handler := New("", serviceShortener)

	logger.Debugln("handler is initialized")

	middleware := middleware.New("", logger, "")

	//------------------------------------------------------------
	// Create mux with chi
	//------------------------------------------------------------
	r := chi.NewRouter()

	r.Use(middleware.Log)
	r.Use(middleware.Auth)
	r.Use(middleware.Encoding)

	//------------------------------------------------------------
	// Use CreateShortURL
	//------------------------------------------------------------
	r.Post("/", handler.CreateShortURL)
}
