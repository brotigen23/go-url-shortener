package app

import (
	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
	"go.uber.org/zap"
)

func Example() {
	var logger *zap.SugaredLogger

	//... Logger init

	var config *config.Config

	//... Config init with godotenv and cleanenv

	err := server.Run(config, logger)
	if err != nil {
		logger.Errorln(err)
	}
}
