package repository

import (
	"context"

	"github.com/okinn/service-presensi/internal/domain/entity"
)

// UserRepository adalah port untuk akses data user
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
