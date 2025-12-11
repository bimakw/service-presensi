package valueobject

type Lokasi struct {
	Latitude  float64
	Longitude float64
	Alamat    string
}

func NewLokasi(lat, long float64, alamat string) *Lokasi {
	if lat == 0 && long == 0 {
		return nil
	}
	return &Lokasi{
		Latitude:  lat,
		Longitude: long,
		Alamat:    alamat,
	}
}
