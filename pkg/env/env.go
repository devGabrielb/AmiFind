package env

import (
	"fmt"
	"os"
)

const (
	PRODUCTION  = "Production"
	DEVELOPMENT = "Development"
)

func TryGetEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok || env == "" {
		err := fmt.Errorf("environment variable not found: %s", key)
		panic(err)
	}
	return env
}

func LoadEnvironment(environmentKey string, fn func(...string) error) {
	v := TryGetEnv(environmentKey)
	if v != PRODUCTION {
		err := fn()
		if err != nil {
			panic(err)
		}
	}
}
