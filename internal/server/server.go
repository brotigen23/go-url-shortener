package server

// main.go

import (
	"net/http"

	"github.com/brotigen23/go-url-shortener/internal/handlers"
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
