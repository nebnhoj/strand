package main

import (
	"log"

	fiber "github.com/gofiber/fiber/v2"

	configs "schuler.com/be-schuler/configs"
	routes "schuler.com/be-schuler/routes"
)

func main() {
	app := fiber.New(configs.SetFiberConfig())
	configs.ConnectDB()
	app.Group("/api").Route("/", routes.BindRoutes)
	log.Fatal(app.Listen(":3000"))

}
