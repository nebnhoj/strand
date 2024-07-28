package todos

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nebnhoj/strand/helpers"
)

func CreateTodo(c *fiber.Ctx) error {
	var todo Todo
	//validate the request body
	if err := c.BodyParser(&todo); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	if errs := helpers.Validator(todo); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}
	newTodo := Todo{
		Id:      uuid.NewString(),
		Name:    todo.Name,
		Details: todo.Details,
	}

	result, err := create(newTodo)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.ResponseSuccess(c, http.StatusOK, result)
}

func FindAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	q := c.Query("q")

	// Convert pagination parameters to integers

	todos, count, err := getAllTodos(q, page, limit)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.ResponsePaginated(c, http.StatusOK, todos, count)
}

type Message struct {
	Content int64 `json:"content"`
}

type Status struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
