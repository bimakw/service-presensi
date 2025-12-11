package entity

import (
	"errors"
	"time"

	"github.com/okinn/service-presensi/internal/domain/valueobject"
)

var (
	ErrInvalidStatus    = errors.New("status presensi tidak valid")
	ErrAlreadyCheckedIn = errors.New("sudah melakukan check-in")
	ErrNotCheckedIn     = errors.New("belum melakukan check-in")
	ErrAlreadyCheckedOut = errors.New("sudah melakukan check-out")
)

type Presensi struct {
	ID         string
	UserID     string
	Nama       string
	Tanggal    time.Time
	JamMasuk   *time.Time
	JamKeluar  *time.Time
	Status     valueobject.StatusPresensi
	Keterangan string
	Lokasi     *valueobject.Lokasi
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewPresensi(userID, nama string, status valueobject.StatusPresensi, keterangan string, lokasi *valueobject.Lokasi) (*Presensi, error) {
	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	now := time.Now()
	p := &Presensi{
		UserID:     userID,
		Nama:       nama,
		Tanggal:    now,
		Status:     status,
		Keterangan: keterangan,
		Lokasi:     lokasi,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Auto set jam masuk jika hadir atau terlambat
	if status == valueobject.StatusHadir || status == valueobject.StatusTerlambat {
		p.JamMasuk = &now
	}

	return p, nil
}

func (p *Presensi) CheckIn() error {
	if p.JamMasuk != nil {
		return ErrAlreadyCheckedIn
	}
	now := time.Now()
	p.JamMasuk = &now
	p.UpdatedAt = now
	return nil
}

func (p *Presensi) CheckOut() error {
	if p.JamMasuk == nil {
		return ErrNotCheckedIn
	}
	if p.JamKeluar != nil {
		return ErrAlreadyCheckedOut
	}
	now := time.Now()
	p.JamKeluar = &now
	p.UpdatedAt = now
	return nil
}

func (p *Presensi) UpdateStatus(status valueobject.StatusPresensi) error {
	if !status.IsValid() {
		return ErrInvalidStatus
	}
	p.Status = status
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Presensi) UpdateKeterangan(keterangan string) {
	p.Keterangan = keterangan
	p.UpdatedAt = time.Now()
}
