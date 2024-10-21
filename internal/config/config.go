package config

import (
	"flag"
	"os"
)

/*
Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
*/

type Config struct {
	ServerAddress string
	BaseURL       string
}

func NewConfig() *Config {
	ret := &Config{}
	flag.StringVar(&ret.ServerAddress, "a", "localhost:8080", "base host")
	flag.StringVar(&ret.BaseURL, "b", "http://localhost:8080", "base host for aliases")
	flag.Parse()
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		ret.ServerAddress = envRunAddr
	}
	if envRunAddr := os.Getenv("BASE_URL"); envRunAddr != "" {
		ret.BaseURL = envRunAddr
	}
	return ret
}
