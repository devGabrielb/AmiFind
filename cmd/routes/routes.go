package routes

import (
	"database/sql"

	"github.com/devGabrielb/AmiFind/internal/pets/handler"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	MapRoutes()
}

type routes struct {
	fb *fiber.App
	rg fiber.Router
	db *sql.DB
}

func New(fb *fiber.App, db *sql.DB) Router {
	return &routes{fb: fb}
}

func (r *routes) MapRoutes() {
	r.rg = r.fb.Group("/api")
	r.rg.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("rodando bem")
	})
	r.buildCustomerRoutes()
	r.buildPostRoutes()
	r.buildPetRoutes()
}

func (r *routes) buildCustomerRoutes() {
	r.rg.Get("/customer", func(c *fiber.Ctx) error {
		return c.SendString("cusomers")
	})
}

func (r *routes) buildPetRoutes() {
	r.rg.Get("/pets", handler.Get)
}

func (r *routes) buildPostRoutes() {
}
