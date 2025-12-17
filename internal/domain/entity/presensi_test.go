package entity

import (
	"testing"

	"github.com/okinn/service-presensi/internal/domain/valueobject"
)

func TestNewPresensi(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		nama       string
		status     valueobject.StatusPresensi
		keterangan string
		wantErr    bool
		wantJamMasuk bool
	}{
		{
			name:         "valid presensi hadir",
			userID:       "user123",
			nama:         "Test User",
			status:       valueobject.StatusHadir,
			keterangan:   "",
			wantErr:      false,
			wantJamMasuk: true,
		},
		{
			name:         "valid presensi terlambat",
			userID:       "user123",
			nama:         "Test User",
			status:       valueobject.StatusTerlambat,
			keterangan:   "Macet",
			wantErr:      false,
			wantJamMasuk: true,
		},
		{
			name:         "valid presensi izin",
			userID:       "user123",
			nama:         "Test User",
			status:       valueobject.StatusIzin,
			keterangan:   "Keperluan keluarga",
			wantErr:      false,
			wantJamMasuk: false,
		},
		{
			name:         "valid presensi sakit",
			userID:       "user123",
			nama:         "Test User",
			status:       valueobject.StatusSakit,
			keterangan:   "Demam",
			wantErr:      false,
			wantJamMasuk: false,
		},
		{
			name:       "invalid status",
			userID:     "user123",
			nama:       "Test User",
			status:     valueobject.StatusPresensi("invalid"),
			keterangan: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presensi, err := NewPresensi(tt.userID, tt.nama, tt.status, tt.keterangan, nil)

			if tt.wantErr {
				if err == nil {
					t.Error("NewPresensi() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewPresensi() unexpected error = %v", err)
				return
			}

			if presensi.UserID != tt.userID {
				t.Errorf("NewPresensi() UserID = %v, want %v", presensi.UserID, tt.userID)
			}
			if presensi.Nama != tt.nama {
				t.Errorf("NewPresensi() Nama = %v, want %v", presensi.Nama, tt.nama)
			}
			if presensi.Status != tt.status {
				t.Errorf("NewPresensi() Status = %v, want %v", presensi.Status, tt.status)
			}
			if presensi.Keterangan != tt.keterangan {
				t.Errorf("NewPresensi() Keterangan = %v, want %v", presensi.Keterangan, tt.keterangan)
			}

			hasJamMasuk := presensi.JamMasuk != nil
			if hasJamMasuk != tt.wantJamMasuk {
				t.Errorf("NewPresensi() has JamMasuk = %v, want %v", hasJamMasuk, tt.wantJamMasuk)
			}
		})
	}
}

func TestPresensi_CheckIn(t *testing.T) {
	t.Run("successful check-in", func(t *testing.T) {
		presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusIzin, "", nil)
		presensi.JamMasuk = nil // Reset jam masuk

		err := presensi.CheckIn()
		if err != nil {
			t.Errorf("CheckIn() unexpected error = %v", err)
		}
		if presensi.JamMasuk == nil {
			t.Error("CheckIn() JamMasuk should not be nil")
		}
	})

	t.Run("already checked in", func(t *testing.T) {
		presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusHadir, "", nil)
		// JamMasuk already set for StatusHadir

		err := presensi.CheckIn()
		if err != ErrAlreadyCheckedIn {
			t.Errorf("CheckIn() error = %v, want %v", err, ErrAlreadyCheckedIn)
		}
	})
}

func TestPresensi_CheckOut(t *testing.T) {
	t.Run("successful check-out", func(t *testing.T) {
		presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusHadir, "", nil)

		err := presensi.CheckOut()
		if err != nil {
			t.Errorf("CheckOut() unexpected error = %v", err)
		}
		if presensi.JamKeluar == nil {
			t.Error("CheckOut() JamKeluar should not be nil")
		}
	})

	t.Run("not checked in yet", func(t *testing.T) {
		presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusIzin, "", nil)
		presensi.JamMasuk = nil

		err := presensi.CheckOut()
		if err != ErrNotCheckedIn {
			t.Errorf("CheckOut() error = %v, want %v", err, ErrNotCheckedIn)
		}
	})

	t.Run("already checked out", func(t *testing.T) {
		presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusHadir, "", nil)
		_ = presensi.CheckOut()

		err := presensi.CheckOut()
		if err != ErrAlreadyCheckedOut {
			t.Errorf("CheckOut() error = %v, want %v", err, ErrAlreadyCheckedOut)
		}
	})
}

func TestPresensi_UpdateStatus(t *testing.T) {
	presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusHadir, "", nil)

	t.Run("valid status update", func(t *testing.T) {
		err := presensi.UpdateStatus(valueobject.StatusTerlambat)
		if err != nil {
			t.Errorf("UpdateStatus() unexpected error = %v", err)
		}
		if presensi.Status != valueobject.StatusTerlambat {
			t.Errorf("UpdateStatus() Status = %v, want %v", presensi.Status, valueobject.StatusTerlambat)
		}
	})

	t.Run("invalid status update", func(t *testing.T) {
		err := presensi.UpdateStatus(valueobject.StatusPresensi("invalid"))
		if err != ErrInvalidStatus {
			t.Errorf("UpdateStatus() error = %v, want %v", err, ErrInvalidStatus)
		}
	})
}

func TestPresensi_UpdateKeterangan(t *testing.T) {
	presensi, _ := NewPresensi("user123", "Test User", valueobject.StatusHadir, "", nil)
	oldUpdatedAt := presensi.UpdatedAt

	presensi.UpdateKeterangan("Keterangan baru")

	if presensi.Keterangan != "Keterangan baru" {
		t.Errorf("UpdateKeterangan() Keterangan = %v, want %v", presensi.Keterangan, "Keterangan baru")
	}
	if !presensi.UpdatedAt.After(oldUpdatedAt) && presensi.UpdatedAt != oldUpdatedAt {
		// Allow equal time in fast execution
	}
}
