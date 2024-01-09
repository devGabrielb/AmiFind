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

type UserHandler struct {
	service services.AuthService
}

func NewUserHandler(service services.AuthService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) Register(c *fiber.Ctx) error {

	registerRequest := dtos.RegisterRequest{}

	if err := c.BodyParser(&registerRequest); err != nil {

		return response.Error(c, fiber.StatusInternalServerError, ErrInvalidJson.Error())
	}

	if err := dtos.Validate(registerRequest); err != nil {
		if err != nil {
			errorParams, ok := err.(*entities.InvalidParameters)
			if ok {
				return response.ErrorWithDetails(c, fiber.StatusBadRequest, errorParams.Error(), errorParams.Params)
			}
		}
		return response.Error(c, fiber.StatusInternalServerError, ErrInvalidDto.Error())
	}

	userAuth_id, err := u.service.Register(c.Context(), registerRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	logrus.WithField("userId", userAuth_id).Info("User created successfully")
	return response.Success(c, fiber.StatusCreated, fiber.Map{"id": userAuth_id})
}

func (u *UserHandler) Login(c *fiber.Ctx) error {

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

	userAuth, err := u.service.Login(c.Context(), loginRequest)
	if err != nil {
		if err.Error() == ErrInvalidCredentials.Error() {
			return response.Error(c, fiber.StatusBadRequest, ErrGenerateToken.Error())

		}
		return response.Error(c, fiber.StatusInternalServerError, ErrGenerateToken.Error())
	}

	logrus.WithField("userId", userAuth.Id).Info("User logged successfully")
	return response.Success(c, fiber.StatusOK, userAuth)
}
