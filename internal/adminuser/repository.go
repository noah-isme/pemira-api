package adminuser

import (
	"context"

	"pemira-api/internal/shared"
)

type Repository interface {
	List(ctx context.Context, filter ListFilter, pag shared.PaginationParams) ([]User, int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, in CreateInput, passwordHash string) (*User, error)
	Update(ctx context.Context, id int64, in UpdateInput) (*User, error)
	ResetPassword(ctx context.Context, id int64, passwordHash string) error
	SetActive(ctx context.Context, id int64, active bool) (*User, error)
	Delete(ctx context.Context, id int64) error
}
