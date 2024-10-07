package config

import (
	"flag"
	"os"
)

/*
Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
*/

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

var Config = config{}

func InitConfig() {
	flag.StringVar(&Config.ServerAddress, "a", "localhost:8080", "base host")
	flag.StringVar(&Config.BaseURL, "b", "http://localhost:8080", "base host for aliases")
	flag.Parse()
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		Config.ServerAddress = envRunAddr
	}
	if envRunAddr := os.Getenv("BASE_URL"); envRunAddr != "" {
		Config.BaseURL = envRunAddr
	}
}
