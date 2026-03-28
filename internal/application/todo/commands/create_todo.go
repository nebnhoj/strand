package commands

import (
	"context"

	"github.com/google/uuid"
	todoDomain "github.com/nebnhoj/strand/internal/domain/todo"
	"github.com/nebnhoj/strand/pkg/cache"
)

// @name CreateTodoCommand
type CreateTodoCommand struct {
	Name    string `json:"name"    validate:"required"`
	Details string `json:"details"`
}

type CreateTodoHandler struct {
	repo  todoDomain.Repository
	cache cache.Cache
}

func NewCreateTodoHandler(repo todoDomain.Repository, c cache.Cache) *CreateTodoHandler {
	return &CreateTodoHandler{repo: repo, cache: c}
}

func (h *CreateTodoHandler) Handle(ctx context.Context, cmd CreateTodoCommand) (string, error) {
	t := todoDomain.Todo{
		ID:      uuid.NewString(),
		Name:    cmd.Name,
		Details: cmd.Details,
	}

	id, err := h.repo.Create(ctx, t)
	if err != nil {
		return "", err
	}

	h.cache.DeleteByPattern(ctx, "todos:*")
	return id, nil
}
