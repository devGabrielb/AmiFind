package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code   int      `json:"code"`
	Msg    string   `json:"mesage"`
	Errors []string `json:"errors,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}
func Error(c *fiber.Ctx, statusCode int, err string) error {
	return c.Status(statusCode).JSON(Response{
		Code: statusCode,
		Msg:  err,
	})
}
func ErrorWithDetails(c *fiber.Ctx, statusCode int, err string, details []string) error {
	return c.Status(statusCode).JSON(Response{
		Code:   statusCode,
		Msg:    err,
		Errors: details,
	})
}
