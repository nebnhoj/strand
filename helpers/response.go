package helpers

import (
	"github.com/gofiber/fiber/v2"
)

func ResponseSuccess(c *fiber.Ctx, status int, result any) error {
	return c.Status(status).JSON(
		Success{
			Status:  status,
			Message: "success",
			Data:    result,
		})

}

func ResponsePaginated(c *fiber.Ctx, status int, result any, total int64) error {
	return c.Status(status).JSON(
		Success{
			Status:  status,
			Message: "success",
			Data:    result,
			Total:   total,
		})

}

func ResponseError(c *fiber.Ctx, status int, error any) error {
	return c.Status(status).JSON(
		Error{
			Status:  status,
			Message: error})

}

type Success struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int64  `json:"total,omitempty" `
}

type Error struct {
	Status  int `json:"status"`
	Message any `json:"message"`
}
