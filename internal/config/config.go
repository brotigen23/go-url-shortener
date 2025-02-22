package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

// Содержит поля, необходимые для инициализации и запуска сервера
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" env-default:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" env-default:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" env-default:"aliases.txt"`
	DatabaseDSN     string `env:"DATABASE_DSN"`

	JWTSecretKey string `env:"SECRET_KEY" env-default:"secret"`
}

// Конструктор Config. Производит чтение переменных окружения и флагов
func NewConfig() (*Config, error) {
	// Read env
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	// Read flags
	a := flag.String("a", "", "server address")
	b := flag.String("b", "", "base host for aliases")
	f := flag.String("f", "", "Path to file with aliases")
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
