package valueobject

import "testing"

func TestNewLokasi(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		long      float64
		alamat    string
		expectNil bool
	}{
		{
			name:      "valid location",
			lat:       -6.2088,
			long:      106.8456,
			alamat:    "Jakarta",
			expectNil: false,
		},
		{
			name:      "valid location with zero alamat",
			lat:       -6.2088,
			long:      106.8456,
			alamat:    "",
			expectNil: false,
		},
		{
			name:      "zero coordinates returns nil",
			lat:       0,
			long:      0,
			alamat:    "Test",
			expectNil: true,
		},
		{
			name:      "zero lat only",
			lat:       0,
			long:      106.8456,
			alamat:    "Test",
			expectNil: false,
		},
		{
			name:      "zero long only",
			lat:       -6.2088,
			long:      0,
			alamat:    "Test",
			expectNil: false,
		},
		{
			name:      "negative coordinates",
			lat:       -33.8688,
			long:      151.2093,
			alamat:    "Sydney",
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lokasi := NewLokasi(tt.lat, tt.long, tt.alamat)

			if tt.expectNil {
				if lokasi != nil {
					t.Error("NewLokasi() expected nil, got non-nil")
				}
				return
			}

			if lokasi == nil {
				t.Error("NewLokasi() expected non-nil, got nil")
				return
			}

			if lokasi.Latitude != tt.lat {
				t.Errorf("NewLokasi() Latitude = %v, want %v", lokasi.Latitude, tt.lat)
			}
			if lokasi.Longitude != tt.long {
				t.Errorf("NewLokasi() Longitude = %v, want %v", lokasi.Longitude, tt.long)
			}
			if lokasi.Alamat != tt.alamat {
				t.Errorf("NewLokasi() Alamat = %v, want %v", lokasi.Alamat, tt.alamat)
			}
		})
	}
}
