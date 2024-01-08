package dtos

import (
	"fmt"
	"strings"

	"github.com/devGabrielb/AmiFind/internal/entities"
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

func Validate(entity interface{}) error {
	value, ok := entity.(struct{})
	if !ok {
		return entities.NewInvalidParams()
	}
	fieldErrors, err := validateEntity(value)
	if err != nil {
		return nil
	}
	if len(fieldErrors) > 0 {
		e := entities.NewInvalidParams()
		e.AddParameters(fieldErrors)
		return e
	}
	return nil
}

func validateEntity(entity struct{}) ([]string, error) {

	fields := make([]string, 0)

	customErrorMessages := map[string]string{
		"required": "is required and cannot be empty",
		"http_url": "must be a valid HTTP URL",
		"email":    "must be a valid email address",
		"max":      "must be at most %s characters long",
		"e164":     "must be a valid phone number",
	}

	err := v.Struct(entity)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}
	}

	for _, validationErr := range err.(validator.ValidationErrors) {
		fieldName := strings.ToLower(validationErr.Field())
		message, found := customErrorMessages[validationErr.Tag()]

		if !found {
			continue
		}

		field := fmt.Sprintf("The field %s %s.\n", fieldName, fmt.Sprintf(message, validationErr.Param()))
		fields = append(fields, field)
	}
	return fields, nil
}
