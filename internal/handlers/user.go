package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/devGabrielb/AmiFind/cmd/api/response"
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
	repo repositories.Repository
}

func NewUserHandler(repo repositories.Repository) *userHandler {
	return &userHandler{repo: repo}
}

func (u *userHandler) Register(c *fiber.Ctx) error {

	userRequest := dtos.RegisterRequest{}

	if err := c.BodyParser(&userRequest); err != nil {

		logrus.WithField("error", err).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrInvalidJson.Error()})
		return response.Error(c, fiber.StatusInternalServerError, ErrInvalidJson.Error())
	}

	pass, err := utils.EncryptPassword(userRequest.Password)

	if err != nil {

		logrus.WithField("error", err).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrEncryptPassword.Error()})
		return response.Error(c, fiber.StatusInternalServerError, ErrEncryptPassword.Error())
	}

	if err := dtos.Validate(userRequest); err != nil {
		logrus.WithField("errors", strings.Split(err.Error(), "\n")).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrInvalidDto.Error()})
		return response.ErrorWithDetails(c, fiber.StatusInternalServerError, ErrInvalidDto.Error(), strings.Split(err.Error(), "\n"))
	}

	user := entities.User{
		Profile_picture: userRequest.Profile_picture,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        string(pass),
		PhoneNumber:     userRequest.PhoneNumber,
		Location:        userRequest.Location,
	}
	id, err := u.repo.Create(c.Context(), user)

	if err != nil {
		logrus.WithField("error", err.Error()).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrCannotCreateUser.Error()})
		return response.Error(c, fiber.StatusInternalServerError, ErrCannotCreateUser.Error())
	}
	logrus.WithField("user_id", id).Info("User created successfully")
	return response.Success(c, fiber.StatusCreated, nil)
}

func (u *userHandler) Login(c *fiber.Ctx) error {

	loginRequest := dtos.LoginRequest{}

	if err := c.BodyParser(&loginRequest); err != nil {

		logrus.WithField("error", err).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrInvalidJson.Error()})
		return response.Error(c, fiber.StatusInternalServerError, ErrInvalidJson.Error())
	}
	if err := dtos.Validate(loginRequest); err != nil {
		logrus.WithField("errors", strings.Split(err.Error(), "\n")).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrInvalidDto.Error()})
		return response.ErrorWithDetails(c, fiber.StatusInternalServerError, ErrInvalidDto.Error(), strings.Split(err.Error(), "\n"))

	}
	user, err := u.repo.FindByEmail(c.Context(), loginRequest.Email)

	if err != nil {

		logrus.WithField("email", loginRequest.Email).Debug(response.Response{Code: fiber.StatusNotFound, Msg: ErrNotFound.Error()})
		return response.Error(c, fiber.StatusNotFound, ErrNotFound.Error())
	}

	logrus.WithField("user_id", user.Id).Info("Found user successfully")

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {

		logrus.WithField("error", err).Debug(response.Response{Code: fiber.StatusBadRequest, Msg: ErrInvalidCredentials.Error()})
		return response.Error(c, fiber.StatusBadRequest, ErrInvalidCredentials.Error())
	}

	env, err := env.TryGetEnv("SECRET_KEY")

	if err != nil {

		logrus.WithField("error", err).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrGetSecretKey.Error()})
		return response.Error(c, fiber.StatusInternalServerError, ErrGetSecretKey.Error())
	}

	t := services.NewToken(env)
	token, err := t.GenerateToken(strconv.Itoa(user.Id))

	if err != nil {

		logrus.WithField("error", err).Error(response.Response{Code: fiber.StatusInternalServerError, Msg: ErrGenerateToken.Error()})
		return response.Error(c, fiber.StatusInternalServerError, ErrGenerateToken.Error())
	}

	logrus.WithField("user_id", user.Id).Info("User logged successfully")
	return response.Success(c, fiber.StatusOK, dtos.LoginResponse{Id: user.Id, Email: user.Email, Token: token})
}
