package router

import (
	"github.com/gofiber/fiber/v3"
	authDomain "github.com/nebnhoj/strand/internal/domain/auth"
	"github.com/nebnhoj/strand/internal/interfaces/http/handlers"
	"github.com/nebnhoj/strand/internal/interfaces/http/middleware"
)

func Bind(
	app fiber.Router,
	jwt *middleware.JWTMiddleware,
	auth *handlers.AuthHandler,
	users *handlers.UserHandler,
	todos *handlers.TodoHandler,
	cacheHandler *handlers.CacheHandler,
) {
	api := app.Group("/api")

	api.Post("/auth", auth.Authenticate)

	todoGroup := api.Group("/todos", jwt.Protected())
	todoGroup.Get("", jwt.RequirePermission(authDomain.PermTodosRead), todos.List)
	todoGroup.Post("", jwt.RequirePermission(authDomain.PermTodosWrite), todos.Create)

	userGroup := api.Group("/users", jwt.Protected())
	userGroup.Get("", jwt.RequirePermission(authDomain.PermUsersRead), users.List)
	userGroup.Get("/:id", jwt.RequirePermission(authDomain.PermUsersRead), users.Get)
	userGroup.Post("", jwt.RequirePermission(authDomain.PermUsersWrite), users.Create)

	api.Delete("/cache", jwt.Protected(), jwt.RequirePermission(authDomain.PermUsersWrite), cacheHandler.Flush)
}
