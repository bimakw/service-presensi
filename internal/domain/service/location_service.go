/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package service

import (
	"context"
	"math"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

const (
	// EarthRadiusMeters is the approximate radius of Earth in meters
	EarthRadiusMeters = 6371000
)

// LocationService handles geofencing logic
type LocationService struct {
	locationRepo repository.AllowedLocationRepository
	enabled      bool
}

// NewLocationService creates a new location service
func NewLocationService(repo repository.AllowedLocationRepository, enabled bool) *LocationService {
	return &LocationService{
		locationRepo: repo,
		enabled:      enabled,
	}
}

// IsEnabled returns whether geofencing is enabled
func (s *LocationService) IsEnabled() bool {
	return s.enabled
}

// ValidateCheckInLocation validates if the given coordinates are within any allowed location
func (s *LocationService) ValidateCheckInLocation(ctx context.Context, lat, lon float64) error {
	if !s.enabled {
		return nil // Geofencing disabled, allow all
	}

	// Skip validation if coordinates are not provided
	if lat == 0 && lon == 0 {
		return nil
	}

	locations, err := s.locationRepo.GetAllActive(ctx)
	if err != nil {
		return err
	}

	if len(locations) == 0 {
		return entity.ErrNoAllowedLocations
	}

	for _, loc := range locations {
		distance := s.HaversineDistance(lat, lon, loc.Latitude, loc.Longitude)
		if distance <= loc.RadiusMeters {
			return nil // Within allowed area
		}
	}

	return entity.ErrOutsideAllowedArea
}

// GetNearestLocation returns the nearest allowed location and distance
func (s *LocationService) GetNearestLocation(ctx context.Context, lat, lon float64) (*entity.AllowedLocation, float64, error) {
	locations, err := s.locationRepo.GetAllActive(ctx)
	if err != nil {
		return nil, 0, err
	}

	if len(locations) == 0 {
		return nil, 0, entity.ErrNoAllowedLocations
	}

	var nearest *entity.AllowedLocation
	minDistance := math.MaxFloat64

	for i := range locations {
		distance := s.HaversineDistance(lat, lon, locations[i].Latitude, locations[i].Longitude)
		if distance < minDistance {
			minDistance = distance
			nearest = &locations[i]
		}
	}

	return nearest, minDistance, nil
}

// HaversineDistance calculates the distance between two coordinates in meters
// using the Haversine formula
func (s *LocationService) HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusMeters * c
}

// IsWithinRadius checks if a point is within the radius of a location
func (s *LocationService) IsWithinRadius(lat, lon float64, location *entity.AllowedLocation) bool {
	distance := s.HaversineDistance(lat, lon, location.Latitude, location.Longitude)
	return distance <= location.RadiusMeters
}
