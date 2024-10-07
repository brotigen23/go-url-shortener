package server

// main.go

import (
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func Run() error {

	// Mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler)

	// Handlers
	//mux.HandleFunc("/status", handlers.StatusHandler)

	// Run server
	return http.ListenAndServe(":8080", mux)
}

func RunWithChi() error {
	r := chi.NewRouter()
	r.Get("/{id}", handlers.IndexGET)
	r.Post("/", handlers.IndexPOST)

	return http.ListenAndServe(config.Config.ServerAddress, r)
}
