package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nebnhoj/strand/middlewares/jwt"
	auth "github.com/nebnhoj/strand/modules/auth"
	schools "github.com/nebnhoj/strand/modules/schools"
	users "github.com/nebnhoj/strand/modules/users"
)

func BindRoutes(router fiber.Router) {
	router.Route("/", UserRoute)
	router.Route("/", SchoolRoute)
	router.Route("/", AuthRoute)

}

func UserRoute(router fiber.Router) {
	userRoute := router.Group("/users", jwt.Protected(), jwt.HasAdminRole)
	userRoute.Get("", users.GetUsers).Name("Get Users")
	userRoute.Get("/:id", users.GetUser).Name("Get User By ID")
	userRoute.Post("", users.CreateUser).Name("Create User")
}
func SchoolRoute(router fiber.Router) {
	router.Get("/schools", schools.ShowHelloWorld).Name("Get Schools")
}

func AuthRoute(router fiber.Router) {
	router.Post("/auth", auth.Authenticate).Name("Authenticate User")
}
