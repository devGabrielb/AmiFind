package handlers

import (
	"errors"

	"github.com/devGabrielb/AmiFind/api/response"
	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("user not found")
)

type UserHandler struct {
	service services.AuthService
}

func NewUserHandler(service services.AuthService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) Register(c *fiber.Ctx) error {

	registerRequest := entities.RegisterRequest{}
	if err := c.BodyParser(&registerRequest); err != nil {
		return response.UnprocessableEntity(c)
	}

	if err := entities.Validate(registerRequest); err != nil {
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

	user := entities.User{
		ProfilePicture: registerRequest.ProfilePicture,
		Name:           registerRequest.Name,
		Email:          registerRequest.Email,
		Password:       registerRequest.Password,
		PhoneNumber:    registerRequest.PhoneNumber,
		Location:       registerRequest.Location,
	}

	userId, err := u.service.Register(c.Context(), user)
	if err != nil {
		if err.Error() == ErrNotFound.Error() {
			return response.NotFound(c)
		}
		return response.InternalServerError(c)
	}

	logrus.WithField("userId", userId).Info("User created successfully")
	return response.Created(c, fiber.Map{"id": userId})
}

func (u *UserHandler) Login(c *fiber.Ctx) error {

	loginRequest := entities.LoginRequest{}
	if err := c.BodyParser(&loginRequest); err != nil {
		return response.UnprocessableEntity(c)
	}

	if err := entities.Validate(loginRequest); err != nil {
		if len(err) > 0 {

			_, ok := err[0].(*entities.ValidateError)
			if ok {
				return response.BadRequestWithErrors(c, err)
			}
			return response.BadRequest(c)
		}
	}

	user, err := u.service.Login(c.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		if err.Error() == ErrNotFound.Error() {
			return response.NotFound(c)
		}
		if err.Error() == ErrInvalidCredentials.Error() {
			return response.BadRequest(c)
		}
		return response.InternalServerError(c)
	}

	token, err := u.service.GenerateToken(user)
	if err != nil {
		return response.InternalServerError(c)

	}

	logrus.WithField("userId", user.Id).Info("User logged successfully")
	return response.OK(c, fiber.Map{"id": user.Id, "token": token})
}
