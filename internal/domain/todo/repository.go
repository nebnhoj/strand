package todo

import "context"

type Repository interface {
	Create(ctx context.Context, t Todo) (string, error)
	FindAll(ctx context.Context, q string, page, limit int) ([]Todo, int64, error)
}
