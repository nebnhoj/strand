package routes

import (
	"github.com/gofiber/fiber/v2"
	schools "schuler.com/be-schuler/modules/schools"
	users "schuler.com/be-schuler/modules/users"
)

func BindRoutes(router fiber.Router) {
	router.Route("/", UserRoute)
	router.Route("/", SchoolRoute)
}

func UserRoute(router fiber.Router) {
	router.Get("/users", users.GetUsers).Name("Get Users")
	router.Get("/users/:id", users.GetUser).Name("Get User By ID")

	router.Post("/users", users.CreateUser).Name("Create User")

}
func SchoolRoute(router fiber.Router) {
	router.Get("/schools", schools.ShowHelloWorld).Name("Get Schools")
}
