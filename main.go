package main

import (
	"log"

	"github.com/gofiber/contrib/swagger"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	configs "github.com/nebnhoj/strand/configs"
	_ "github.com/nebnhoj/strand/docs"
	routes "github.com/nebnhoj/strand/routes"
)

// @title Your API
// @version 1.0
// @description This is a sample server for a to-do application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/
func main() {
	app := fiber.New(configs.SetFiberConfig())
	cfg := swagger.Config{
		BasePath: "/api/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
		Title:    "Swagger API Docs",
	}

	app.Use(swagger.New(cfg))
	app.Use(logger.New())
	routes.BindRoutes(app)

	log.Fatal(app.Listen(":3000"))

}
