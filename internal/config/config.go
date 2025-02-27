package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

/*
	"server_address": "localhost:8080",
    "base_url": "http://localhost",
    "file_storage_path": "/path/to/file.db",
    "database_dsn": "",
    "enable_https"
*/
// Содержит поля, необходимые для инициализации и запуска сервера
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" env-default:"localhost:8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" env-default:"http://localhost:8080" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" env-default:"aliases.txt" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`

	JWTSecretKey string `env:"SECRET_KEY" env-default:"secret"`

	EnableHTTPS bool `env:"ENABLE_HTTPS" env-default:"false" json:"enable_https"`

	ConfigFile string `env:"CONFIG"`
}

func readJSONConfig(path string) (*Config, error) {

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Конструктор Config. Производит чтение переменных окружения и флагов
func NewConfig() (*Config, error) {
	// Read env
	cfg := &Config{}
	var err error
	// Flags
	a := flag.String("a", "", "server address")
	b := flag.String("b", "", "base host for aliases")
	f := flag.String("f", "", "Path to file with aliases")
	d := flag.String("d", "", "String connection to DB")
	c := flag.String("c", "", "Config file")
	s := flag.Bool("s", false, "Enable HTTPS")

	flag.Parse()
	log.Println("BEFORE")
	// JSON
	if *c != "" {
		log.Println("IN")

		cfg, err = readJSONConfig(*c)
		if err != nil {
			return nil, err
		}
	}

	// ENV
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
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
	cfg.EnableHTTPS = *s
	return cfg, nil
}
