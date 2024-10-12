package app

import (
	"fmt"
	"os"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
	"github.com/brotigen23/go-url-shortener/internal/storage"
)

func Run() {
	config := config.NewConfig()
	storage := storage.NewStorage()


	err := server.Run(config, storage)



	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
