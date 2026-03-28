package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	todoCmds "github.com/nebnhoj/strand/internal/application/todo/commands"
	todoQrs "github.com/nebnhoj/strand/internal/application/todo/queries"
	"github.com/nebnhoj/strand/pkg/response"
	"github.com/nebnhoj/strand/pkg/validator"
)

type TodoHandler struct {
	create    *todoCmds.CreateTodoHandler
	listTodos *todoQrs.ListTodosHandler
}

func NewTodoHandler(
	create *todoCmds.CreateTodoHandler,
	listTodos *todoQrs.ListTodosHandler,
) *TodoHandler {
	return &TodoHandler{create: create, listTodos: listTodos}
}

// Create godoc
// @Summary Create todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param body body CreateTodoRequest true "New todo payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /todos [post]
func (h *TodoHandler) Create(c fiber.Ctx) error {
	var cmd todoCmds.CreateTodoCommand
	if err := c.Bind().Body(&cmd); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	if errs := validator.Validate(cmd); len(errs) > 0 {
		return response.Error(c, http.StatusBadRequest, errs)
	}

	id, err := h.create.Handle(c.Context(), cmd)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusCreated, map[string]string{"id": id})
}

// List godoc
// @Summary List todos
// @Description Paginated list of todos
// @Tags todos
// @Produce json
// @Param page  query int    false "Page number"
// @Param limit query int    false "Page size"
// @Param q     query string false "Search query"
// @Success 200 {array}  TodoResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /todos [get]
func (h *TodoHandler) List(c fiber.Ctx) error {
	todos, total, err := h.listTodos.Handle(c.Context(), todoQrs.ListTodosQuery{
		Q:     c.Query("q"),
		Page:  atoi(c.Query("page")),
		Limit: atoi(c.Query("limit")),
	})
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Paginated(c, http.StatusOK, todos, total)
}
