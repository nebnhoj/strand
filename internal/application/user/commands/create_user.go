package commands

import (
	"context"
	"strings"

	"github.com/google/uuid"
	userDomain "github.com/nebnhoj/strand/internal/domain/user"
	"github.com/nebnhoj/strand/pkg/cache"
	"golang.org/x/crypto/bcrypt"
)

// @name CreateUserCommand
type CreateUserCommand struct {
	FirstName       string             `json:"first_name"       validate:"required"`
	LastName        string             `json:"last_name"        validate:"required"`
	Title           string             `json:"title"`
	Email           string             `json:"email"            validate:"required,email"`
	Password        string             `json:"password"         validate:"required"`
	ConfirmPassword string             `json:"confirm_password" validate:"required,eqfield=Password"`
	Roles           []string           `json:"roles"            validate:"required"`
	Address         userDomain.Address `json:"address"          validate:"required"`
}

type CreateUserHandler struct {
	repo  userDomain.Repository
	cache cache.Cache
}

func NewCreateUserHandler(repo userDomain.Repository, c cache.Cache) *CreateUserHandler {
	return &CreateUserHandler{repo: repo, cache: c}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u := userDomain.User{
		ID:        uuid.NewString(),
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Title:     cmd.Title,
		Email:     strings.ToLower(cmd.Email),
		Password:  string(hashed),
		Roles:     cmd.Roles,
		Address:   cmd.Address,
	}

	id, err := h.repo.Create(ctx, u)
	if err != nil {
		return "", err
	}

	h.cache.DeleteByPattern(ctx, "users:*")
	return id, nil
}
