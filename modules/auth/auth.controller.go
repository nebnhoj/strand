package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nebnhoj/strand/helpers"
)

func Authenticate(c *fiber.Ctx) error {
	var auth Auth
	if err := c.BodyParser(&auth); err != nil {
		return helpers.ResponseError(c, http.StatusBadRequest, err)
	}

	if errs := helpers.Validator(auth); len(errs) > 0 && errs[0].Error {
		return helpers.ResponseError(c, http.StatusBadRequest, errs)
	}

	token, err := GetJWTToken(auth)
	if err != nil {
		return helpers.ResponseError(c, http.StatusUnauthorized, err.Error())
	}

	return helpers.ResponseSuccess(c, http.StatusOK, token)
}
