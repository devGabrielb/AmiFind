package handlers

import (
	"log"
	"strconv"

	"github.com/devGabrielb/AmiFind/internal/services"
	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service services.PostService
}

func NewPostHandler(service services.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (p *PostHandler) Get(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return err
	}

	posts, err := p.service.Get(c.Context(), page)
	log.Println(posts.Data)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(posts.Data)
}
