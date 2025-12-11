package usecase

import (
	"context"
	"errors"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
	"github.com/okinn/service-presensi/internal/domain/valueobject"
)

var (
	ErrPresensiNotFound = errors.New("presensi tidak ditemukan")
)

// CreatePresensiInput adalah input untuk membuat presensi
type CreatePresensiInput struct {
	UserID     string
	Nama       string
	Status     string
	Keterangan string
	Latitude   float64
	Longitude  float64
	Alamat     string
}

// UpdatePresensiInput adalah input untuk update presensi
type UpdatePresensiInput struct {
	Status     string
	Keterangan string
}

// PresensiOutput adalah output untuk presensi
type PresensiOutput struct {
	ID         string                     `json:"id"`
	UserID     string                     `json:"user_id"`
	Nama       string                     `json:"nama"`
	Tanggal    string                     `json:"tanggal"`
	JamMasuk   *string                    `json:"jam_masuk,omitempty"`
	JamKeluar  *string                    `json:"jam_keluar,omitempty"`
	Status     valueobject.StatusPresensi `json:"status"`
	Keterangan string                     `json:"keterangan,omitempty"`
	Lokasi     *LokasiOutput              `json:"lokasi,omitempty"`
	CreatedAt  string                     `json:"created_at"`
	UpdatedAt  string                     `json:"updated_at"`
}

type LokasiOutput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Alamat    string  `json:"alamat,omitempty"`
}

// PresensiUseCase adalah interface untuk use case presensi
type PresensiUseCase interface {
	Create(ctx context.Context, input CreatePresensiInput) (*PresensiOutput, error)
	GetByID(ctx context.Context, id string) (*PresensiOutput, error)
	GetAll(ctx context.Context, filter repository.PresensiFilter, page, limit int) ([]PresensiOutput, int64, error)
	Update(ctx context.Context, id string, input UpdatePresensiInput) (*PresensiOutput, error)
	Delete(ctx context.Context, id string) error
	CheckIn(ctx context.Context, id string) error
	CheckOut(ctx context.Context, id string) error
}

type presensiUseCase struct {
	repo repository.PresensiRepository
}

func NewPresensiUseCase(repo repository.PresensiRepository) PresensiUseCase {
	return &presensiUseCase{repo: repo}
}

func (uc *presensiUseCase) Create(ctx context.Context, input CreatePresensiInput) (*PresensiOutput, error) {
	lokasi := valueobject.NewLokasi(input.Latitude, input.Longitude, input.Alamat)

	presensi, err := entity.NewPresensi(
		input.UserID,
		input.Nama,
		valueobject.StatusPresensi(input.Status),
		input.Keterangan,
		lokasi,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, presensi); err != nil {
		return nil, err
	}

	return toPresensiOutput(presensi), nil
}

func (uc *presensiUseCase) GetByID(ctx context.Context, id string) (*PresensiOutput, error) {
	presensi, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPresensiNotFound
	}
	return toPresensiOutput(presensi), nil
}

func (uc *presensiUseCase) GetAll(ctx context.Context, filter repository.PresensiFilter, page, limit int) ([]PresensiOutput, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	presensiList, total, err := uc.repo.GetAll(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]PresensiOutput, len(presensiList))
	for i, p := range presensiList {
		outputs[i] = *toPresensiOutput(&p)
	}

	return outputs, total, nil
}

func (uc *presensiUseCase) Update(ctx context.Context, id string, input UpdatePresensiInput) (*PresensiOutput, error) {
	presensi, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPresensiNotFound
	}

	if input.Status != "" {
		if err := presensi.UpdateStatus(valueobject.StatusPresensi(input.Status)); err != nil {
			return nil, err
		}
	}
	if input.Keterangan != "" {
		presensi.UpdateKeterangan(input.Keterangan)
	}

	if err := uc.repo.Update(ctx, presensi); err != nil {
		return nil, err
	}

	return toPresensiOutput(presensi), nil
}

func (uc *presensiUseCase) Delete(ctx context.Context, id string) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrPresensiNotFound
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *presensiUseCase) CheckIn(ctx context.Context, id string) error {
	presensi, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrPresensiNotFound
	}

	if err := presensi.CheckIn(); err != nil {
		return err
	}

	return uc.repo.Update(ctx, presensi)
}

func (uc *presensiUseCase) CheckOut(ctx context.Context, id string) error {
	presensi, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return ErrPresensiNotFound
	}

	if err := presensi.CheckOut(); err != nil {
		return err
	}

	return uc.repo.Update(ctx, presensi)
}

func toPresensiOutput(p *entity.Presensi) *PresensiOutput {
	output := &PresensiOutput{
		ID:         p.ID,
		UserID:     p.UserID,
		Nama:       p.Nama,
		Tanggal:    p.Tanggal.Format("2006-01-02"),
		Status:     p.Status,
		Keterangan: p.Keterangan,
		CreatedAt:  p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if p.JamMasuk != nil {
		jamMasuk := p.JamMasuk.Format("15:04:05")
		output.JamMasuk = &jamMasuk
	}
	if p.JamKeluar != nil {
		jamKeluar := p.JamKeluar.Format("15:04:05")
		output.JamKeluar = &jamKeluar
	}
	if p.Lokasi != nil {
		output.Lokasi = &LokasiOutput{
			Latitude:  p.Lokasi.Latitude,
			Longitude: p.Lokasi.Longitude,
			Alamat:    p.Lokasi.Alamat,
		}
	}

	return output
}
