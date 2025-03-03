package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Log      Log      `json:"log" envPrefix:"LOG_"`
	Server   Server   `json:"server" envPrefix:"SERVER_"`
	Postgres Postgres `json:"postgres" envPrefix:"PG_"`

	OpenapiPath string `json:"openapi_path" env:"OPENAPI_PATH"`
	AppVersion  string `json:"app_version" env:"-"`
}

type Log struct {
	Level       int    `json:"log_level" env:"LEVEL"`
	VictoriaUrl string `json:"victoria_url" env:"VICTORIA_URL"`
}

type Server struct {
	Port         int           `json:"port" env:"PORT"`
	ReadTimeout  time.Duration `json:"read_timeout" env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `json:"write_timeout" env:"WRITE_TIMEOUT" envDefault:"5s"`
}

type Postgres struct {
	RoDSN          string        `json:"ro_dsn" env:"RO_DSN"`
	RwDSN          string        `json:"rw_dsn" env:"RW_DSN"`
	ConnectTimeout time.Duration `json:"connect_timeout" env:"CONNECT_TIMEOUT" envDefault:"5s"`
}

func GetDefault() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to parse: %w", err)
	}

	return &cfg, nil
}
