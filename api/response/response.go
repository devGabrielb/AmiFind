package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrBadRequest          = errors.New("the server could not understand the request due to invalid syntax or missing parameters")
	ErrForbiden            = errors.New("the server understood the request, but it refuses to authorize it. The client does not have the necessary permissions")
	ErrInternalServerError = errors.New("an unexpected error occurred on the server. Please try again later")
	ErrUnprocessableEntity = errors.New("the server understands the request, but it cannot process the provided data due to semantic errors")
	ErrNotFound            = errors.New("the requested resource could not be found on the server")
)

type Response struct {
	Code   int      `json:"code"`
	Msg    string   `json:"mesage"`
	Errors []string `json:"errors,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}
func Error(c *fiber.Ctx, statusCode int, err error) error {
	return c.Status(statusCode).JSON(Response{
		Code: statusCode,
		Msg:  err.Error(),
	})
}
func ErrorWithDetails(c *fiber.Ctx, statusCode int, errorMesage error, errs []error) error {
	errArray := []string{}
	for _, err := range errs {
		errArray = append(errArray, err.Error())
	}
	return c.Status(statusCode).JSON(Response{
		Code:   statusCode,
		Msg:    errorMesage.Error(),
		Errors: errArray,
	})
}

func OK(c *fiber.Ctx, data interface{}) error {
	return Success(c, fiber.StatusOK, data)
}

func Created(c *fiber.Ctx, data interface{}) error {
	return Success(c, fiber.StatusCreated, data)
}

func BadRequest(c *fiber.Ctx) error {
	return Error(c, fiber.StatusBadRequest, ErrBadRequest)
}

func BadRequestWithErrors(c *fiber.Ctx, innerErrors []error) error {
	return ErrorWithDetails(c, fiber.StatusBadRequest, ErrBadRequest, innerErrors)
}

func Forbiden(c *fiber.Ctx) error {
	return Error(c, fiber.StatusForbidden, ErrForbiden)
}
func UnprocessableEntity(c *fiber.Ctx) error {
	return Error(c, fiber.StatusUnprocessableEntity, ErrUnprocessableEntity)
}

func InternalServerError(c *fiber.Ctx) error {
	return Error(c, fiber.StatusInternalServerError, ErrInternalServerError)
}
func NotFound(c *fiber.Ctx) error {
	return Error(c, fiber.StatusNotFound, ErrNotFound)
}
