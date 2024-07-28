package routes

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/nebnhoj/strand/docs"
	middlewares "github.com/nebnhoj/strand/middlewares/jwt"
	auth "github.com/nebnhoj/strand/modules/auth"
	todos "github.com/nebnhoj/strand/modules/todos"
	users "github.com/nebnhoj/strand/modules/users"
)

// BindRoutes sets up all the routes for the application.
// It registers the user, todo, and authentication routes to the provided router.
func BindRoutes(router fiber.Router) {

	api := router.Group("/api")

	// Register the routes for user-related operations
	api.Route("/", UserRoutes)
	// Register the routes for todo-related operations
	api.Route("/", TodoRoutes)
	// Register the routes for authentication-related operations
	api.Route("/", AuthRoutes)
}

func UserRoutes(router fiber.Router) {
	userRoute := router.Group("/users", middlewares.Protected(), middlewares.HasAdminRole)

	// @Summary Get Users
	// @Description Get all users
	// @Tags users
	// @Produce json
	// @Success 200 {array} users.User
	// @Router /users [get]
	userRoute.Get("", users.GetUsers).Name("Get Users")

	// @Summary Get User By ID
	// @Description Get a user by ID
	// @Tags users
	// @Produce json
	// @Param id path int true "User ID"
	// @Success 200 {object} users.User
	// @Router /users/{id} [get]
	userRoute.Get("/:id", users.GetUser).Name("Get User By ID")

	// @Summary Create User
	// @Description Create a new user
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param user body users.User true "New User"
	// @Success 201 {object} users.User
	// @Router /users [post]
	userRoute.Post("", users.CreateUser).Name("Create User")
}

func TodoRoutes(router fiber.Router) {

	// @Summary Create Todo
	// @Description Create a new todo
	// @Tags todos
	// @Accept json
	// @Produce json
	// @Param todo body todos.Todo true "New Todo"
	// @Success 201 {object} todos.Todo
	// @Router /todos [post]
	router.Post("/todos", todos.CreateTodo).Name("Create Todos")

	// @Summary Get Todos
	// @Description Get all todos
	// @Tags todos
	// @Produce json
	// @Success 200 {array} todos.Todo
	// @Router /todos [get]
	router.Get("/todos", todos.FindAll).Name("Get Todos")
}

func AuthRoutes(router fiber.Router) {

	// @Summary Authenticate User
	// @Description Authenticate a user
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param credentials body auth.Credentials true "User Credentials"
	// @Success 200 {object} auth.Token
	// @Router /auth [post]
	router.Post("/auth", auth.Authenticate).Name("Authenticate User")
}
