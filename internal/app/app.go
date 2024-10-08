package app

import (
	"fmt"

	"github.com/brotigen23/go-url-shortener/internal/config"
	"github.com/brotigen23/go-url-shortener/internal/server"
)

func Fun() {
	//config.InitConfig()
	config.InitConfigENV()
	fmt.Println(config.ConfigENV.Host)
	fmt.Println(config.ConfigENV.BastHostForAliases)
	server.RunWithChi()
}
