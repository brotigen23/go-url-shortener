package app

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
)

func Run() error {
	config := config.NewConfig()
	err := server.Run(config)

	if err != nil {
		return err
	}
	return nil
}
