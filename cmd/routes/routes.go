package routes

import (
	"database/sql"

	"github.com/devGabrielb/AmiFind/internal/handlers"
	"github.com/devGabrielb/AmiFind/internal/repositories"
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
	return &routes{fb: fb, db: db}
}

func (r *routes) MapRoutes() {
	r.rg = r.fb.Group("/api")
	r.buildUserRoutes()
	r.buildPostRoutes()
	r.buildPetRoutes()
}

func (r *routes) buildUserRoutes() {
	repo := repositories.NewRepository(r.db)
	handle := handlers.NewUserHandler(repo)

	r.rg.Post("/register", handle.Register)
	r.rg.Post("/register", handle.Login)
}

func (r *routes) buildPetRoutes() {
}

func (r *routes) buildPostRoutes() {
}
