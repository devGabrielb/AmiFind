package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	PRODUCTION  = "Production"
	DEVELOPMENT = "Development"
)

func TryGetEnv(key string) (string, error) {
	env, ok := os.LookupEnv(key)
	if !ok || env == "" {
		return "", fmt.Errorf("environment variable not found: %s", key)
	}
	return env, nil
}

func TryGetEnvList(keys ...string) (mapKey map[string]string, errorList error) {
	mapKey = make(map[string]string)
	for _, k := range keys {
		key, err := TryGetEnv(k)
		if err != nil {
			errorList = errors.Join(err)
		}
		mapKey[k] = key

		if errorList != nil {
			return map[string]string{}, errorList
		}

	}

	return mapKey, nil
}
