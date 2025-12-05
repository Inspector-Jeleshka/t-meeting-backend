package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort int    `env:"HTTP_PORT" env-default:"33"`
	DBDSN    string `env:"DB_DSN" env-default:"postgres://postgres:coolpassword@localhost:5433/tmeeting?sslmode=disable"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}
	return &cfg
}
