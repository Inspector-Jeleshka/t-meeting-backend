package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type DatabaseConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5433"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"tmeeting"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
}

type Config struct {
	Env      string         `yaml:"env"       env:"APP_ENV"   env-default:"local"`
	HTTPPort int            `yaml:"http_port" env:"HTTP_PORT" env-default:"33"`
	DB       DatabaseConfig `yaml:"db"`
}

// сделал метод, чтобы был формат строки, подходящий под pgx
func (c *Config) DBDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}
	return &cfg
}
