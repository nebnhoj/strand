package main

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	configs "github.com/nebnhoj/strand/configs"
	routes "github.com/nebnhoj/strand/routes"
)

func main() {
	app := fiber.New(configs.SetFiberConfig())
	configs.ConnectDB()
	app.Group("/api").Route("/", routes.BindRoutes)
	app.Use(logger.New())
	log.Fatal(app.Listen(":3000"))

}
