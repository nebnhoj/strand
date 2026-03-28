package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"

	_ "github.com/nebnhoj/strand/docs"

	authCmds "github.com/nebnhoj/strand/internal/application/auth/commands"
	todoCmds "github.com/nebnhoj/strand/internal/application/todo/commands"
	todoQrs "github.com/nebnhoj/strand/internal/application/todo/queries"
	userCmds "github.com/nebnhoj/strand/internal/application/user/commands"
	userQrs "github.com/nebnhoj/strand/internal/application/user/queries"

	infrajwt "github.com/nebnhoj/strand/internal/infrastructure/jwt"
	"github.com/nebnhoj/strand/internal/infrastructure/mongodb"
	infraredis "github.com/nebnhoj/strand/internal/infrastructure/redis"

	"github.com/nebnhoj/strand/internal/interfaces/http/handlers"
	"github.com/nebnhoj/strand/internal/interfaces/http/middleware"
	"github.com/nebnhoj/strand/internal/interfaces/router"
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
	godotenv.Load()

	// ── Infrastructure ────────────────────────────────────────────────────────
	db := mongodb.ConnectDB()
	jwtSecret := os.Getenv("JWT_SECRET")
	redisClient := infraredis.NewClient()
	appCache := infraredis.NewCache(redisClient)

	// ── Repositories (domain interfaces, mongo implementations) ──────────────
	userRepo := mongodb.NewUserRepository(db, "strand")
	todoRepo := mongodb.NewTodoRepository(db, "strand")

	// ── Services ──────────────────────────────────────────────────────────────
	tokenSvc := infrajwt.NewTokenService(jwtSecret)

	// ── Command handlers ──────────────────────────────────────────────────────
	authenticateCmd := authCmds.NewAuthenticateHandler(userRepo, tokenSvc)
	createUserCmd := userCmds.NewCreateUserHandler(userRepo, appCache)
	createTodoCmd := todoCmds.NewCreateTodoHandler(todoRepo, appCache)

	// ── Query handlers ────────────────────────────────────────────────────────
	listUsersQ := userQrs.NewListUsersHandler(userRepo, appCache)
	getUserQ := userQrs.NewGetUserHandler(userRepo, appCache)
	listTodosQ := todoQrs.NewListTodosHandler(todoRepo, appCache)

	// ── HTTP handlers ─────────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authenticateCmd)
	userHandler := handlers.NewUserHandler(createUserCmd, getUserQ, listUsersQ)
	todoHandler := handlers.NewTodoHandler(createTodoCmd, listTodosQ)
	cacheHandler := handlers.NewCacheHandler(appCache)

	// ── Middleware ────────────────────────────────────────────────────────────
	jwtMw := middleware.NewJWTMiddleware(tokenSvc)

	// ── Fiber app ─────────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  os.Getenv("APP_HEADER"),
		AppName:       os.Getenv("APP_NAME"),
	})

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
	router.Bind(app, jwtMw, authHandler, userHandler, todoHandler, cacheHandler)

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
