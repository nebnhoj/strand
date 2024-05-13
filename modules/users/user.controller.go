package users

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	helpers "github.com/nebnhoj/strand/helpers"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *fiber.Ctx) error {

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user := claims["user"]
	log.Printf("User: %s", user)

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
	var user UserDTO
	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	if errs := helpers.Validator(user); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}

	newUser := User{
		Id:        uuid.NewString(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Title:     user.Title,
		Email:     strings.ToLower(user.Email),
		Password:  Hash(user.Password),
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

type Message struct {
	Content string `json:"content"`
}

func Hash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hashed := string(hashedPassword)
	log.Printf("Hashed Password: %s", hashed)
	return hashed
}
