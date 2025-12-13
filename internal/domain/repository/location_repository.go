/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package repository

import (
	"context"

	"github.com/okinn/service-presensi/internal/domain/entity"
)

// AllowedLocationRepository defines the interface for allowed location persistence
type AllowedLocationRepository interface {
	// Create creates a new allowed location
	Create(ctx context.Context, location *entity.AllowedLocation) error

	// GetByID retrieves an allowed location by ID
	GetByID(ctx context.Context, id string) (*entity.AllowedLocation, error)

	// GetAll retrieves all allowed locations
	GetAll(ctx context.Context) ([]entity.AllowedLocation, error)

	// GetAllActive retrieves all active allowed locations
	GetAllActive(ctx context.Context) ([]entity.AllowedLocation, error)

	// Update updates an allowed location
	Update(ctx context.Context, location *entity.AllowedLocation) error

	// Delete deletes an allowed location by ID
	Delete(ctx context.Context, id string) error
}
