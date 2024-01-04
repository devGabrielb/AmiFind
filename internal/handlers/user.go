package handlers

import (
	"errors"
	"strconv"

	"github.com/devGabrielb/AmiFind/internal/dtos"
	"github.com/devGabrielb/AmiFind/internal/repositories"
	"github.com/devGabrielb/AmiFind/internal/utils"

	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/services"
	"github.com/devGabrielb/AmiFind/pkg/env"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidJson        = errors.New("invalid json")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidDto         = errors.New("invalid request")
	ErrCannotCreateUser   = errors.New("cannot create user")
	ErrEncryptPassword    = errors.New("error while encrypting password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrGenerateToken      = errors.New("error generating token")
	ErrGetSecretKey       = errors.New("error getting secret key")
)

type userHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *userHandler {
	return &userHandler{repo: repo}
}

func (u *userHandler) Register(c *fiber.Ctx) error {

	userRequest := dtos.RegisterRequest{}

	if err := c.BodyParser(&userRequest); err != nil {

		logrus.WithField("error", err).Error(ErrInvalidJson.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.NewError(fiber.StatusInternalServerError, ErrInvalidJson.Error()))
	}

	pass, err := utils.EncryptPassword(userRequest.Password)

	if err != nil {

		logrus.WithField("error", err).Error(ErrEncryptPassword.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.NewError(fiber.StatusInternalServerError, ErrEncryptPassword.Error()))
	}

	user := entities.User{
		Profile_picture: userRequest.Profile_picture_url,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        string(pass),
		PhoneNumber:     userRequest.PhoneNumber,
		Location:        userRequest.Location,
	}
	id, err := u.repo.Create(c.Context(), user)

	if err != nil {
		logrus.WithField("error", err.Error()).Error(ErrCannotCreateUser.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.NewError(fiber.StatusInternalServerError, ErrCannotCreateUser.Error()))
	}
	logrus.WithField("user_id", id).Info("User created successfully")
	return c.Status(201).JSON(nil)
}

func (u *userHandler) Login(c *fiber.Ctx) error {

	loginRequest := dtos.LoginRequest{}

	if err := c.BodyParser(&loginRequest); err != nil {

		logrus.WithField("error", err).Error(ErrInvalidJson.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.NewError(fiber.StatusInternalServerError, ErrInvalidJson.Error()))
	}

	user, err := u.repo.FindByEmail(c.Context(), loginRequest.Email)

	if err != nil {

		logrus.WithField("email", loginRequest.Email).Debug(dtos.NewError(fiber.StatusNotFound, ErrNotFound.Error()))
		return c.Status(fiber.StatusNotFound).JSON(dtos.NewError(fiber.StatusNotFound, ErrNotFound.Error()))
	}

	logrus.WithField("user_id", user.Id).Info("Found user successfully")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {

		logrus.WithField("error", err).Debug(dtos.NewError(fiber.StatusBadRequest, ErrInvalidCredentials.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(dtos.NewError(fiber.StatusBadRequest, ErrInvalidCredentials.Error()))
	}

	env, err := env.TryGetEnv("SECRET_KEY")

	if err != nil {

		logrus.WithField("error", err).Error(dtos.NewError(fiber.StatusInternalServerError, ErrGetSecretKey.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.NewError(fiber.StatusInternalServerError, ErrGetSecretKey.Error()))
	}

	t := services.NewToken(env)
	token, err := t.GenerateToken(strconv.Itoa(user.Id))

	if err != nil {

		logrus.WithField("error", err).Error(dtos.NewError(fiber.StatusInternalServerError, ErrGenerateToken.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(dtos.LoginResponse{Id: user.Id, Email: user.Email, Token: token})
}
