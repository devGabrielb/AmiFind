package dtos

import (
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

type RegisterRequest struct {
	Profile_picture string `json:"profile_picture" validate:"required,http_url"`
	Name            string `json:"name" validate:"required,max=24"`
	Email           string `json:"email" validate:"email,required,max=24"`
	Password        string `json:"password" validate:"required,max=12"`
	PhoneNumber     string `json:"phoneNumber" validate:"e164,required,max=20"`
	Location        string `json:"location" validate:"required,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,max=24"`
	Password string `json:"password" validate:"required,max=12"`
}

func (u *LoginRequest) Validate() error {
	return nil
}

type LoginResponse struct {
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
	Id    int    `json:"id,omitempty"`
}

func Validate(model interface{}) error {

	v = validator.New(validator.WithRequiredStructEnabled())
	err := v.Struct(model)
	if err != nil {
		return err

	}
	return nil
}
