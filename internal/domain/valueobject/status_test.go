package valueobject

import "testing"

func TestStatusPresensi_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   StatusPresensi
		expected bool
	}{
		{"hadir is valid", StatusHadir, true},
		{"terlambat is valid", StatusTerlambat, true},
		{"izin is valid", StatusIzin, true},
		{"sakit is valid", StatusSakit, true},
		{"alpha is valid", StatusAlpha, true},
		{"invalid status", StatusPresensi("invalid"), false},
		{"empty status", StatusPresensi(""), false},
		{"random string", StatusPresensi("random"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.expected {
				t.Errorf("StatusPresensi.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatusPresensi_String(t *testing.T) {
	tests := []struct {
		name     string
		status   StatusPresensi
		expected string
	}{
		{"hadir string", StatusHadir, "hadir"},
		{"terlambat string", StatusTerlambat, "terlambat"},
		{"izin string", StatusIzin, "izin"},
		{"sakit string", StatusSakit, "sakit"},
		{"alpha string", StatusAlpha, "alpha"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("StatusPresensi.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
