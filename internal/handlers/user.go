package handlers

import (
	"strconv"

	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/repositories"
	"github.com/devGabrielb/AmiFind/internal/services"
	"github.com/devGabrielb/AmiFind/pkg/env"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterRequest struct {
	Profile_picture_url string `json:"profilePictureUrl" validate:"required,max=255"`
	Name                string `json:"name" validate:"required,max=24"`
	Email               string `json:"email" validate:"required,max=24"`
	Password            string `json:"password" validate:"required,max=12"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,max=20"`
	Location            string `json:"location" validate:"required,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=24"`
	Password string `json:"password" validate:"required,max=12"`
}

type LoginResponse struct {
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
	Id    int    `json:"id,omitempty"`
}

func (u *UserRegisterRequest) Validate() error {
	return nil
}

type userHandler struct {
	repo repositories.UserRepository
}

func NewUserHandler(repo repositories.UserRepository) *userHandler {
	return &userHandler{repo: repo}
}

func (u *userHandler) Register(c *fiber.Ctx) error {
	userRequest := UserRegisterRequest{}

	if err := c.BodyParser(&userRequest); err != nil {
		return err
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 5)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Message": "Erro ao criptografar senha"})
	}

	user := entities.User{
		Profile_picture: userRequest.Profile_picture_url,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        string(pass),
		PhoneNumber:     userRequest.PhoneNumber,
		Location:        userRequest.Location,
	}
	if err := u.repo.Create(c.Context(), user); err != nil {
		logrus.WithField("error", err.Error()).Error("error while request user by email")
		return err
	}

	return c.Status(201).JSON(nil)
}

func (u *userHandler) Login(c *fiber.Ctx) error {
	loginRequest := LoginRequest{}
	if err := c.BodyParser(&loginRequest); err != nil {
		return err
	}

	user, err := u.repo.FindByEmail(c.Context(), loginRequest.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	env, err := env.TryGetEnv("SECRET_KEY")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	t := services.NewToken(env)

	token, err := t.GenerateToken(strconv.Itoa(user.Id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(LoginResponse{Id: user.Id, Email: user.Email, Token: token})
}
