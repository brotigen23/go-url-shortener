package app

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
)

func Fun() {
	config.InitConfig()
	server.RunWithChi()
}
