package auth

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/nebnhoj/strand/helpers"
)

// Authenticate godoc
// @Summary Authenticate user
// @Description Authenticate with email and password, returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Auth true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} helpers.Error
// @Failure 401 {object} helpers.Error
// @Router /auth [post]
func Authenticate(c fiber.Ctx) error {
	var auth Auth
	if err := c.Bind().Body(&auth); err != nil {
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
