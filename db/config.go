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

func tryGetConfigDB() (*Config, error) {
	envs, err := env.TryGetEnvList("DB_HOST", "DB_PORT", "DB_USER", "DB_NAME", "DB_PASS")
	if err != nil {
		return nil, err
	}
	return &Config{
		database: struct {
			Host   string
			Port   string
			User   string
			DbName string
			Pass   string
		}{
			Host:   envs["DB_HOST"],
			Port:   envs["DB_PORT"],
			User:   envs["DB_USER"],
			DbName: envs["DB_NAME"],
			Pass:   envs["DB_PASS"],
		},
	}, nil
}
