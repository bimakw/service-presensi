package repository

import (
	"context"
	"time"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/valueobject"
)

// PresensiFilter untuk filtering data presensi
type PresensiFilter struct {
	UserID    string
	Status    valueobject.StatusPresensi
	StartDate time.Time
	EndDate   time.Time
}

// PresensiRepository adalah port untuk akses data presensi
// Interface ini didefinisikan di domain, implementasi di adapter
type PresensiRepository interface {
	Create(ctx context.Context, presensi *entity.Presensi) error
	GetByID(ctx context.Context, id string) (*entity.Presensi, error)
	GetAll(ctx context.Context, filter PresensiFilter, page, limit int) ([]entity.Presensi, int64, error)
	Update(ctx context.Context, presensi *entity.Presensi) error
	Delete(ctx context.Context, id string) error
}
