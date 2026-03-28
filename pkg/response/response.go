package response

import "github.com/gofiber/fiber/v3"

// ErrorBody is exported so swagger can reference it in annotations.
// @name ErrorBody
type ErrorBody struct {
	Status  int `json:"status"`
	Message any `json:"message"`
}

type successBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int64  `json:"total,omitempty"`
}

func Success(c fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(successBody{Status: status, Message: "success", Data: data})
}

func Paginated(c fiber.Ctx, status int, data any, total int64) error {
	return c.Status(status).JSON(successBody{Status: status, Message: "success", Data: data, Total: total})
}

func Error(c fiber.Ctx, status int, msg any) error {
	return c.Status(status).JSON(ErrorBody{Status: status, Message: msg})
}
