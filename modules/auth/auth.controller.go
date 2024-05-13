package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nebnhoj/strand/helpers"
	appJWT "github.com/nebnhoj/strand/middlewares/jwt"
	"golang.org/x/crypto/bcrypt"

	errors "github.com/nebnhoj/strand/helpers/errors"
	users "github.com/nebnhoj/strand/modules/users"
)

func Authenticate(c *fiber.Ctx) error {
	var auth Auth
	//validate the request body
	if err := c.BodyParser(&auth); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	if errs := helpers.Validator(auth); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}
	user, err := users.GetUserByEmail(auth.Email)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {

		return helpers.ResponseError(c, http.StatusUnauthorized, errors.UNAUTHORIZE)
	}

	t, err := appJWT.CreateJWTClaim(user)
	if err != nil {
		return helpers.ResponseError(c, http.StatusInternalServerError, err.Error())
	}
	return helpers.ResponseSuccess(c, http.StatusOK, t)
}
