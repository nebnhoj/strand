package routes

import (
	"github.com/gofiber/fiber/v2"
	middlewares "github.com/nebnhoj/strand/middlewares/jwt"
	auth "github.com/nebnhoj/strand/modules/auth"
	todos "github.com/nebnhoj/strand/modules/todos"
	users "github.com/nebnhoj/strand/modules/users"
)

func BindRoutes(router fiber.Router) {
	router.Route("/", UserRoutes)
	router.Route("/", TodoRoutes)
	router.Route("/", AuthRoutes)

}

func UserRoutes(router fiber.Router) {
	userRoute := router.Group("/users", middlewares.Protected(), middlewares.HasAdminRole)
	userRoute.Get("", users.GetUsers).Name("Get Users")
	userRoute.Get("/:id", users.GetUser).Name("Get User By ID")
	userRoute.Post("", users.CreateUser).Name("Create User")
}
func TodoRoutes(router fiber.Router) {
	router.Post("/todos", todos.CreateTodo).Name("Create Todos")
	router.Get("/todos", todos.FindAll).Name("Get Todos")

}

func AuthRoutes(router fiber.Router) {
	router.Post("/auth", auth.Authenticate).Name("Authenticate User")
}
