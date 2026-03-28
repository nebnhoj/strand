package todos

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/nebnhoj/strand/helpers"
)

// CreateTodo godoc
// @Summary Create a todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body Todo true "New todo payload"
// @Success 200 {object} Todo
// @Failure 400 {object} helpers.Error
// @Failure 500 {object} helpers.Error
// @Router /todos [post]
func CreateTodo(c fiber.Ctx) error {
	var todo Todo
	if err := c.Bind().Body(&todo); err != nil {
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

// FindAll godoc
// @Summary Get all todos
// @Description Get paginated list of todos
// @Tags todos
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param q query string false "Search query"
// @Success 200 {array} Todo
// @Failure 500 {object} helpers.Error
// @Router /todos [get]
func FindAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	q := c.Query("q")

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
