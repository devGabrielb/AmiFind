package db

import (
	"github.com/devGabrielb/AmiFind/pkg/env"
)

type Config struct {
	database struct {
		Host   string
		Port   string
		User   string
		DbName string
		Pass   string
	}
}

func newConfig() *Config {
	return &Config{
		database: struct {
			Host   string
			Port   string
			User   string
			DbName string
			Pass   string
		}{
			Host:   env.TryGetEnv("DB_HOST"),
			Port:   env.TryGetEnv("DB_PORT"),
			User:   env.TryGetEnv("DB_USER"),
			DbName: env.TryGetEnv("DB_NAME"),
			Pass:   env.TryGetEnv("DB_PASS"),
		},
	}
}
