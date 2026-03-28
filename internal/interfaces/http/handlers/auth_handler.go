package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	authCmds "github.com/nebnhoj/strand/internal/application/auth/commands"
	"github.com/nebnhoj/strand/pkg/response"
	"github.com/nebnhoj/strand/pkg/validator"
)

type AuthHandler struct {
	authenticate *authCmds.AuthenticateHandler
}

func NewAuthHandler(authenticate *authCmds.AuthenticateHandler) *AuthHandler {
	return &AuthHandler{authenticate: authenticate}
}

// Authenticate godoc
// @Summary Authenticate user
// @Description Returns a JWT token for valid credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param body body AuthRequest true "Credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth [post]
func (h *AuthHandler) Authenticate(c fiber.Ctx) error {
	var cmd authCmds.AuthenticateCommand
	if err := c.Bind().Body(&cmd); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	if errs := validator.Validate(cmd); len(errs) > 0 {
		return response.Error(c, http.StatusBadRequest, errs)
	}

	result, err := h.authenticate.Handle(c.Context(), cmd)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, err.Error())
	}

	return response.Success(c, http.StatusOK, result)
}
