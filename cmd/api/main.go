package main

import (
	"log"
	"os"

	"github.com/devGabrielb/AmiFind/cmd/routes"
	"github.com/devGabrielb/AmiFind/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var err error

func main() {
	env := os.Getenv("ENV_KEY")

	if env != "Production" {
		err = godotenv.Load()
	}
	if err != nil {
		log.Printf(".env not loadedA")
	}
	app := fiber.New()

	db, err := db.NewDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	routes := routes.New(app, db)
	routes.MapRoutes()
	log.Fatal(app.Listen(":9090"))
}
