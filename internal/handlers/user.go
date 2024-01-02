package handlers

import (
	"github.com/devGabrielb/AmiFind/internal/entities"
	"github.com/devGabrielb/AmiFind/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type UserRegisterRequest struct {
	Profile_picture_url string `json:"profilePictureUrl" validate:"required,max=255"`
	Name                string `json:"name" validate:"required,max=24"`
	Email               string `json:"email" validate:"required,max=24"`
	Password            string `json:"password" validate:"required,max=12"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,max=20"`
	Location            string `json:"location" validate:"required,max=255"`
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
	user := entities.User{
		Profile_picture: userRequest.Profile_picture_url,
		Name:            userRequest.Name,
		Email:           userRequest.Email,
		Password:        userRequest.Password,
		PhoneNumber:     userRequest.PhoneNumber,
		Location:        userRequest.Location,
	}
	err := u.repo.Create(c.Context(), user)
	if err != nil {
		return err
	}

	return nil
}
