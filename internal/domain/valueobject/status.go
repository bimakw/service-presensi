package valueobject

type StatusPresensi string

const (
	StatusHadir     StatusPresensi = "hadir"
	StatusTerlambat StatusPresensi = "terlambat"
	StatusIzin      StatusPresensi = "izin"
	StatusSakit     StatusPresensi = "sakit"
	StatusAlpha     StatusPresensi = "alpha"
)

func (s StatusPresensi) IsValid() bool {
	switch s {
	case StatusHadir, StatusTerlambat, StatusIzin, StatusSakit, StatusAlpha:
		return true
	}
	return false
}

func (s StatusPresensi) String() string {
	return string(s)
}
