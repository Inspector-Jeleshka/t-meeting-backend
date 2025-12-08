package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env"       env:"APP_ENV"   env-default:"local"`
	HTTPPort int    `yaml:"http_port" env:"HTTP_PORT" env-default:"33"`
	DBDSN    string `yaml:"-"    env:"DB_DSN"    env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}
	return &cfg
}
