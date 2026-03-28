package queries

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	todoDomain "github.com/nebnhoj/strand/internal/domain/todo"
	"github.com/nebnhoj/strand/pkg/cache"
)

type ListTodosQuery struct {
	Q     string
	Page  int
	Limit int
}

type todoListResult struct {
	Items []todoDomain.Todo `json:"items"`
	Total int64             `json:"total"`
}

type ListTodosHandler struct {
	repo  todoDomain.Repository
	cache cache.Cache
}

func NewListTodosHandler(repo todoDomain.Repository, c cache.Cache) *ListTodosHandler {
	return &ListTodosHandler{repo: repo, cache: c}
}

func (h *ListTodosHandler) Handle(ctx context.Context, q ListTodosQuery) ([]todoDomain.Todo, int64, error) {
	key := fmt.Sprintf("todos:list:%s:%d:%d", q.Q, q.Page, q.Limit)

	if data, err := h.cache.Get(ctx, key); err == nil {
		var cached todoListResult
		if err := json.Unmarshal(data, &cached); err == nil {
			return cached.Items, cached.Total, nil
		}
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		_ = err // cache unavailable — fall through to DB
	}

	items, total, err := h.repo.FindAll(ctx, q.Q, q.Page, q.Limit)
	if err != nil {
		return nil, 0, err
	}

	if data, merr := json.Marshal(todoListResult{Items: items, Total: total}); merr == nil {
		h.cache.Set(ctx, key, data, 5*time.Minute)
	}
	return items, total, nil
}
