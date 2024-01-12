package entities

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var customErrorMessages = map[string]string{
	"required": "is required and cannot be empty",
	"http_url": "must be a valid HTTP URL",
	"email":    "must be a valid email address",
	"max":      "must be at most %s characters long",
	"e164":     "must be a valid phone number",
}

type ValidateError struct {
	Msg string `json:"message"`
}

func NewValidateError(msg string) *ValidateError {
	return &ValidateError{Msg: msg}
}
func (i *ValidateError) Error() string {
	return i.Msg
}

func Validate(entity interface{}) []error {
	v := validator.New(validator.WithRequiredStructEnabled())

	fields := make([]error, 0)

	err := v.Struct(entity)

	if err != nil {
		if e, ok := err.(*validator.InvalidValidationError); ok {

			return []error{e}
		}

		for _, validationErr := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(validationErr.Field())
			message, found := customErrorMessages[validationErr.Tag()]

			if !found {
				continue
			}

			field := NewValidateError(fmt.Sprintf("The field %s %s.", fieldName, fmt.Sprintf(message, validationErr.Param())))
			fields = append(fields, field)
		}
	}

	return fields
}
