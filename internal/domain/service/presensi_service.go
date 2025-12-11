package service

import (
	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/valueobject"
)

// PresensiDomainService berisi business logic yang tidak cocok di entity
type PresensiDomainService struct{}

func NewPresensiDomainService() *PresensiDomainService {
	return &PresensiDomainService{}
}

// DetermineStatus menentukan status berdasarkan jam masuk
func (s *PresensiDomainService) DetermineStatus(jamMasuk int, batasJam int) valueobject.StatusPresensi {
	if jamMasuk <= batasJam {
		return valueobject.StatusHadir
	}
	return valueobject.StatusTerlambat
}

// CalculateDuration menghitung durasi kerja
func (s *PresensiDomainService) CalculateDuration(presensi *entity.Presensi) float64 {
	if presensi.JamMasuk == nil || presensi.JamKeluar == nil {
		return 0
	}
	duration := presensi.JamKeluar.Sub(*presensi.JamMasuk)
	return duration.Hours()
}
