package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	userCmds "github.com/nebnhoj/strand/internal/application/user/commands"
	userQrs "github.com/nebnhoj/strand/internal/application/user/queries"
	"github.com/nebnhoj/strand/pkg/response"
	"github.com/nebnhoj/strand/pkg/validator"
)

type UserHandler struct {
	create    *userCmds.CreateUserHandler
	getUser   *userQrs.GetUserHandler
	listUsers *userQrs.ListUsersHandler
}

func NewUserHandler(
	create *userCmds.CreateUserHandler,
	getUser *userQrs.GetUserHandler,
	listUsers *userQrs.ListUsersHandler,
) *UserHandler {
	return &UserHandler{create: create, getUser: getUser, listUsers: listUsers}
}

// List godoc
// @Summary List users
// @Description Paginated list of users (admin only)
// @Tags users
// @Produce json
// @Param page  query int    false "Page number"
// @Param limit query int    false "Page size"
// @Param q     query string false "Search query"
// @Success 200 {array}  UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	users, err := h.listUsers.Handle(c.Context(), userQrs.ListUsersQuery{
		Q:     c.Query("q"),
		Page:  atoi(c.Query("page")),
		Limit: atoi(c.Query("limit")),
	})
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}
	return response.Success(c, http.StatusOK, users)
}

// Get godoc
// @Summary Get user by ID
// @Description Get a single user by their ID (admin only)
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) Get(c fiber.Ctx) error {
	u, err := h.getUser.Handle(c.Context(), userQrs.GetUserQuery{ID: c.Params("id")})
	if err != nil {
		return response.Error(c, http.StatusNotFound, err.Error())
	}
	return response.Success(c, http.StatusOK, u)
}

// Create godoc
// @Summary Create user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param body body CreateUserRequest true "New user payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var cmd userCmds.CreateUserCommand
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

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
