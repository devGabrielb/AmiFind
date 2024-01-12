package handlers

import (
	"strconv"

	"github.com/devGabrielb/AmiFind/api/response"
	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AdvertisingHandler struct {
	service services.AdvertisingService
}

func NewAdvertisingHandler(service services.AdvertisingService) *AdvertisingHandler {
	return &AdvertisingHandler{service: service}
}

func (a *AdvertisingHandler) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		response.Forbiden(c)
	}
	Id := claims["userId"].(string)

	userId, err := strconv.Atoi(Id)
	if err != nil {
		return err
	}
	advertisingrequest := entities.AdvertisingRequest{}
	if err := c.BodyParser(&advertisingrequest); err != nil {
		return response.UnprocessableEntity(c)
	}

	if err := entities.Validate(advertisingrequest); err != nil {
		if len(err) > 0 {

			if err != nil {
				_, ok := err[0].(*entities.ValidateError)
				if ok {
					return response.BadRequestWithErrors(c, err)
				}
			}
			return response.BadRequest(c)
		}
	}
	advertising := &entities.Advertising{
		Status:   advertisingrequest.Status,
		Category: advertisingrequest.Category,
	}

	advertising.UserId = userId
	advertising.NewPost(advertisingrequest.Post)
	advertising.Post.NewPet(advertisingrequest.Post.Pet)

	id, err := a.service.CreateAd(c.Context(), *advertising)
	if err != nil {
		return response.InternalServerError(c)
	}
	return response.Created(c, fiber.Map{"id": id})
}

func (a *AdvertisingHandler) GetByQuery(c *fiber.Ctx) error {
	queries := c.Queries()
	_, err := a.service.GetByQueyParams(c.Context(), queries)
	if err != nil {
		return response.NotFound(c)
	}
	return err
}
