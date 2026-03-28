package user

import "context"

type Repository interface {
	Create(ctx context.Context, u User) (string, error)
	FindAll(ctx context.Context, q string, page, limit int) ([]User, error)
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}
