package users

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	helpers "schuler.com/be-schuler/helpers"
)

func GetUsers(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	q := c.Query("q")

	// Convert pagination parameters to integers
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	users, err := getAllUsers(q, page, limit)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}

	return helpers.ResponseSuccess(c, http.StatusOK, users)
}

func GetUser(c *fiber.Ctx) error {
	user, err := getUserByID(c.Params("id"))
	if err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	return helpers.ResponseSuccess(c, http.StatusOK, user)
}

func CreateUser(c *fiber.Ctx) error {
	var user User
	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	//use the validator library to validate required fields

	if errs := helpers.Validator(user); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}

	newUser := User{
		Id:       uuid.NewString(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
		Email:    user.Email,
	}

	result, err := create(newUser)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err)
	}

	return helpers.ResponseSuccess(c, http.StatusCreated, result)

}

type Message struct {
	Content string `json:"content"`
}
