package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

/*
Флаг -a отвечает за адрес запуска HTTP-сервера (значение может быть таким: localhost:8888).
Флаг -b отвечает за базовый адрес результирующего сокращённого URL (значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
*/

type config struct {
	BaseHost           *string
	BastHostForAliases *string
}

var Config = config{}

func InitConfig() {
	Config.BaseHost = flag.String("a", "localhost:8080", "base host")
	Config.BastHostForAliases = flag.String("b", "http://localhost:8080", "base host for aliases")
	flag.Parse()
}

type configENV struct {
	Host               string `env:"SERVERBASEHOST"`
	BastHostForAliases string `env:"BASEHOSTFORALIASES"`
}

var ConfigENV = configENV{}

func InitConfigENV() {
	err := env.Parse(&ConfigENV)
	if err != nil {
		log.Fatal(err)
	}
}
