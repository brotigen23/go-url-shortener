package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" env-default:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" env-default:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`

	JWTSecretKey string `env:"SECRET_KEY" env-default:"secret"`
}

func NewConfig() (*Config, error) {
	// Read env
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	// Read flags
	a := flag.String("a", "localhost:8080", "server address")
	b := flag.String("b", "http://localhost:8080", "base host for aliases")
	f := flag.String("f", "./aliases.txt", "Path to file with aliases")
	d := flag.String("d", "", "String connection to DB")

	flag.Parse()
	if *a != "" {
		cfg.ServerAddress = *a
	}
	if *b != "" {
		cfg.BaseURL = *b
	}
	if *d != "" {
		cfg.DatabaseDSN = *d
	}
	if *f != "" {
		cfg.FileStoragePath = *f
	}
	return cfg, nil
}
