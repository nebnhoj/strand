package commands

import (
	"context"
	"errors"

	authDomain "github.com/nebnhoj/strand/internal/domain/auth"
	userDomain "github.com/nebnhoj/strand/internal/domain/user"
	"github.com/nebnhoj/strand/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

// @name AuthenticateCommand
type AuthenticateCommand struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResult struct {
	Token       string   `json:"token"`
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
}

type AuthenticateHandler struct {
	users  userDomain.Repository
	tokens authDomain.TokenService
}

func NewAuthenticateHandler(users userDomain.Repository, tokens authDomain.TokenService) *AuthenticateHandler {
	return &AuthenticateHandler{users: users, tokens: tokens}
}

func (h *AuthenticateHandler) Handle(ctx context.Context, cmd AuthenticateCommand) (AuthResult, error) {
	u, err := h.users.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return AuthResult{}, errors.New(apperrors.UNAUTHORIZED)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(cmd.Password)); err != nil {
		return AuthResult{}, errors.New(apperrors.UNAUTHORIZED)
	}

	token, err := h.tokens.Generate(u.ID, u.Email, u.Roles, u.Permissions)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, UserID: u.ID, Email: u.Email, Permissions: u.Permissions}, nil
}
