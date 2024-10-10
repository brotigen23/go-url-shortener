package app

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
)

func Run() {
	config := config.NewConfig()
	server.Run(config)
}
