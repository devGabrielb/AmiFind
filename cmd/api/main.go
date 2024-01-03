package main

import (
	"log"
	"os"

	"github.com/devGabrielb/AmiFind/cmd/routes"
	"github.com/devGabrielb/AmiFind/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV_KEY")

	if env != "Production" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
	app := fiber.New()

	db, err := db.NewDb()
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	routes := routes.New(app, db)
	routes.MapRoutes()

	log.Fatal(app.Listen(":9090"))
}
