/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidLocationName   = errors.New("nama lokasi tidak valid")
	ErrInvalidCoordinates    = errors.New("koordinat tidak valid")
	ErrInvalidRadius         = errors.New("radius harus lebih dari 0")
	ErrOutsideAllowedArea    = errors.New("lokasi check-in di luar area yang diizinkan")
	ErrNoAllowedLocations    = errors.New("tidak ada lokasi yang dikonfigurasi")
	ErrGeofencingDisabled    = errors.New("geofencing tidak aktif")
)

// AllowedLocation represents a location where check-in is permitted
type AllowedLocation struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`          // e.g., "Kantor Pusat", "Cabang Jakarta"
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RadiusMeters float64   `json:"radius_meters"` // Allowed check-in radius in meters
	Address      string    `json:"address"`       // Human-readable address
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewAllowedLocation creates a new allowed location with validation
func NewAllowedLocation(name string, lat, lon, radiusMeters float64, address string) (*AllowedLocation, error) {
	if name == "" {
		return nil, ErrInvalidLocationName
	}

	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return nil, ErrInvalidCoordinates
	}

	if radiusMeters <= 0 {
		return nil, ErrInvalidRadius
	}

	now := time.Now()
	return &AllowedLocation{
		Name:         name,
		Latitude:     lat,
		Longitude:    lon,
		RadiusMeters: radiusMeters,
		Address:      address,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update updates the allowed location
func (l *AllowedLocation) Update(name string, lat, lon, radiusMeters float64, address string, isActive bool) error {
	if name == "" {
		return ErrInvalidLocationName
	}

	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return ErrInvalidCoordinates
	}

	if radiusMeters <= 0 {
		return ErrInvalidRadius
	}

	l.Name = name
	l.Latitude = lat
	l.Longitude = lon
	l.RadiusMeters = radiusMeters
	l.Address = address
	l.IsActive = isActive
	l.UpdatedAt = time.Now()

	return nil
}

// Deactivate deactivates the location
func (l *AllowedLocation) Deactivate() {
	l.IsActive = false
	l.UpdatedAt = time.Now()
}

// Activate activates the location
func (l *AllowedLocation) Activate() {
	l.IsActive = true
	l.UpdatedAt = time.Now()
}
