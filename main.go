package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	configs "github.com/nebnhoj/strand/configs"
	_ "github.com/nebnhoj/strand/docs"
	routes "github.com/nebnhoj/strand/routes"
)

// @title Strand API
// @version 1.0
// @description REST API for the Strand application.
// @host localhost:3001
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter: Bearer {token}
func main() {

	app := fiber.New(configs.SetFiberConfig())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}))

	app.Get("/api/docs/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	app.Get("/api/docs", func(c fiber.Ctx) error {
		return c.Type("html").SendString(swaggerUI("/api/docs/swagger.json", "Swagger API Docs"))
	})

	app.Use(logger.New())
	routes.BindRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))

}

func swaggerUI(specURL, title string) string {
	return `<!DOCTYPE html>
<html>
<head>
  <title>` + title + `</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script>
  SwaggerUIBundle({
    url: "` + specURL + `",
    dom_id: '#swagger-ui',
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
    layout: "BaseLayout"
  })
</script>
</body>
</html>`
}
