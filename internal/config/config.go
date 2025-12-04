package config

import (
	"fmt"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type HTTPConfig struct {
	Port int
}

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func MustLoad() *Config {
	cfg := &Config{}

	// http
	httpPortStr := getEnv("HTTP_PORT", "33")
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		panic(fmt.Sprintf("invalid HHTP_PORT: %v", err))
	}
	cfg.HTTP.Port = httpPort

	//db
	dbPortStr := getEnv("DB_PORT", "5433")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		panic(fmt.Sprintf("invalid DB_PORT: %v", err))
	}

	cfg.DB = DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "coolpassword"),
		Name:     getEnv("DB_NAME", "tmeeting"),
	}

	return cfg
}

func (c *Config) DBDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}
