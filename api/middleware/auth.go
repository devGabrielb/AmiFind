package middleware

import (
	"fmt"

	"github.com/devGabrielb/AmiFind/pkg/env"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	env, err := env.TryGetEnv("SECRET_KEY")
	if err != nil {
		fmt.Print(err)
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(env)},
	})
}
