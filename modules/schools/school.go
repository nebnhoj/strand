package schools

import (
	"github.com/gofiber/fiber/v2"
)

func ShowHelloWorld(c *fiber.Ctx) error {
	m := new(Message)
	m.Content = "schools endpoint"
	return c.JSON(m)
}

type Message struct {
	Content string `json:"content"`
}
