package users

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	helpers "github.com/nebnhoj/strand/helpers"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get paginated list of users (admin only)
// @Tags users
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param q query string false "Search query"
// @Success 200 {array} User
// @Failure 401 {object} helpers.Error
// @Failure 403 {object} helpers.Error
// @Failure 500 {object} helpers.Error
// @Security BearerAuth
// @Router /users [get]
func GetUsers(c fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	q := c.Query("q")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	users, err := getAllUsers(q, page, limit)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.ResponseSuccess(c, http.StatusOK, users)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a single user by their ID (admin only)
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} helpers.Error
// @Failure 401 {object} helpers.Error
// @Failure 403 {object} helpers.Error
// @Security BearerAuth
// @Router /users/{id} [get]
func GetUser(c fiber.Ctx) error {
	user, err := getUserByID(c.Params("id"))
	if err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	return helpers.ResponseSuccess(c, http.StatusOK, user)
}

// CreateUser godoc
// @Summary Create a user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserDTO true "New user payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} helpers.Error
// @Failure 401 {object} helpers.Error
// @Failure 403 {object} helpers.Error
// @Failure 500 {object} helpers.Error
// @Security BearerAuth
// @Router /users [post]
func CreateUser(c fiber.Ctx) error {
	var user UserDTO
	if err := c.Bind().Body(&user); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	if errs := helpers.Validator(user); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}

	hashed, err := Hash(user.Password)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err)
	}

	newUser := User{
		Id:        uuid.NewString(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Title:     user.Title,
		Email:     strings.ToLower(user.Email),
		Password:  hashed,
		Roles:     user.Roles,
		Address: Address{
			City:     user.Address.City,
			Street:   user.Address.Street,
			Province: user.Address.Province,
			Country:  user.Address.Country,
		},
	}

	result, err := create(newUser)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err)
	}

	return helpers.ResponseSuccess(c, http.StatusCreated, result)
}

func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
