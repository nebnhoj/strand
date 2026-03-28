package queries

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	userDomain "github.com/nebnhoj/strand/internal/domain/user"
	"github.com/nebnhoj/strand/pkg/cache"
)

type GetUserQuery struct {
	ID string
}

type GetUserHandler struct {
	repo  userDomain.Repository
	cache cache.Cache
}

func NewGetUserHandler(repo userDomain.Repository, c cache.Cache) *GetUserHandler {
	return &GetUserHandler{repo: repo, cache: c}
}

func (h *GetUserHandler) Handle(ctx context.Context, q GetUserQuery) (userDomain.User, error) {
	key := fmt.Sprintf("users:get:%s", q.ID)

	if data, err := h.cache.Get(ctx, key); err == nil {
		var result userDomain.User
		if err := json.Unmarshal(data, &result); err == nil {
			return result, nil
		}
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err // cache unavailable — fall through to DB
	}

	result, err := h.repo.FindByID(ctx, q.ID)
	if err != nil {
		return userDomain.User{}, err
	}

	if data, merr := json.Marshal(result); merr == nil {
		h.cache.Set(ctx, key, data, 5*time.Minute)
	}
	return result, nil
}
