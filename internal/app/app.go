package app

import (
	"fmt"
	"os"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
)

func Run() {
	config := config.NewConfig()
	err := server.Run(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
