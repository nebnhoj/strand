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

type ListUsersQuery struct {
	Q     string
	Page  int
	Limit int
}

type ListUsersHandler struct {
	repo  userDomain.Repository
	cache cache.Cache
}

func NewListUsersHandler(repo userDomain.Repository, c cache.Cache) *ListUsersHandler {
	return &ListUsersHandler{repo: repo, cache: c}
}

func (h *ListUsersHandler) Handle(ctx context.Context, q ListUsersQuery) ([]userDomain.User, error) {
	key := fmt.Sprintf("users:list:%s:%d:%d", q.Q, q.Page, q.Limit)

	if data, err := h.cache.Get(ctx, key); err == nil {
		var result []userDomain.User
		if err := json.Unmarshal(data, &result); err == nil {
			return result, nil
		}
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err // cache unavailable — fall through to DB
	}

	result, err := h.repo.FindAll(ctx, q.Q, q.Page, q.Limit)
	if err != nil {
		return nil, err
	}

	if data, merr := json.Marshal(result); merr == nil {
		h.cache.Set(ctx, key, data, 5*time.Minute)
	}
	return result, nil
}
