package handlers

import (
	"errors"

	"github.com/devGabrielb/AmiFind/cmd/api/response"
	"github.com/devGabrielb/AmiFind/internal/dtos"
	"github.com/devGabrielb/AmiFind/internal/entities"

	"github.com/devGabrielb/AmiFind/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidJson        = errors.New("invalid json")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrGenerateToken      = errors.New("error generating token")
	ErrGetSecretKey       = errors.New("error getting secret key")
	ErrInvalidParams      = errors.New("invalid parameters")
	ErrInvalidDto         = errors.New("invalid request")
)

type userHandler struct {
	service services.Service
}

func NewUserHandler(service services.Service) *userHandler {
	return &userHandler{service: service}
}

func (u *userHandler) Register(c *fiber.Ctx) error {

	registerRequest := dtos.RegisterRequest{}

	if err := c.BodyParser(&registerRequest); err != nil {

		return response.Error(c, fiber.StatusInternalServerError, ErrInvalidJson.Error())
	}

	if err := dtos.Validate(registerRequest); err != nil {
		errorParams, ok := err.(*entities.InvalidParameters)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, ErrInvalidDto.Error())
		}
		return response.ErrorWithDetails(c, fiber.StatusBadRequest, errorParams.Error(), errorParams.Params)
	}

	user_id, err := u.service.Register(c.Context(), registerRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	logrus.WithField("user_id", user_id).Info("User created successfully")
	return response.Success(c, fiber.StatusCreated, fiber.Map{"id": user_id})
}

func (u *userHandler) Login(c *fiber.Ctx) error {

	loginRequest := dtos.LoginRequest{}

	if err := c.BodyParser(&loginRequest); err != nil {

		return response.Error(c, fiber.StatusInternalServerError, ErrInvalidJson.Error())
	}
	if err := dtos.Validate(loginRequest); err != nil {
		errorParams, ok := err.(*entities.InvalidParameters)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, ErrInvalidDto.Error())
		}
		return response.ErrorWithDetails(c, fiber.StatusInternalServerError, errorParams.Error(), errorParams.Params)
	}

	user, err := u.service.Login(c.Context(), loginRequest)
	if err != nil {
		if err.Error() == ErrInvalidCredentials.Error() {
			return response.Error(c, fiber.StatusBadRequest, ErrGenerateToken.Error())

		}
		return response.Error(c, fiber.StatusInternalServerError, ErrGenerateToken.Error())
	}

	logrus.WithField("user_id", user.Id).Info("User logged successfully")
	return response.Success(c, fiber.StatusOK, user)
}
