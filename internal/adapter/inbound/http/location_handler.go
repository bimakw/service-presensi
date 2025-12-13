/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package http

import (
	"encoding/json"
	"net/http"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
	"github.com/okinn/service-presensi/pkg/validator"
)

type LocationHandler struct {
	repo repository.AllowedLocationRepository
}

func NewLocationHandler(repo repository.AllowedLocationRepository) *LocationHandler {
	return &LocationHandler{repo: repo}
}

type CreateLocationRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Latitude     float64 `json:"latitude" validate:"required,gte=-90,lte=90"`
	Longitude    float64 `json:"longitude" validate:"required,gte=-180,lte=180"`
	RadiusMeters float64 `json:"radius_meters" validate:"required,gt=0"`
	Address      string  `json:"address" validate:"max=255"`
}

type UpdateLocationRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Latitude     float64 `json:"latitude" validate:"required,gte=-90,lte=90"`
	Longitude    float64 `json:"longitude" validate:"required,gte=-180,lte=180"`
	RadiusMeters float64 `json:"radius_meters" validate:"required,gt=0"`
	Address      string  `json:"address" validate:"max=255"`
	IsActive     bool    `json:"is_active"`
}

type LocationOutput struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	RadiusMeters float64 `json:"radius_meters"`
	Address      string  `json:"address,omitempty"`
	IsActive     bool    `json:"is_active"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func (h *LocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			ValidationError(w, validationErrs)
			return
		}
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	location, err := entity.NewAllowedLocation(
		req.Name,
		req.Latitude,
		req.Longitude,
		req.RadiusMeters,
		req.Address,
	)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Create(r.Context(), location); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusCreated, "Lokasi berhasil ditambahkan", toLocationOutput(location))
}

func (h *LocationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	locations, err := h.repo.GetAll(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	outputs := make([]LocationOutput, len(locations))
	for i, loc := range locations {
		outputs[i] = *toLocationOutput(&loc)
	}

	Success(w, http.StatusOK, "Berhasil", outputs)
}

func (h *LocationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	location, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, "Lokasi tidak ditemukan")
		return
	}

	Success(w, http.StatusOK, "Berhasil", toLocationOutput(location))
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			ValidationError(w, validationErrs)
			return
		}
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	location, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, "Lokasi tidak ditemukan")
		return
	}

	if err := location.Update(
		req.Name,
		req.Latitude,
		req.Longitude,
		req.RadiusMeters,
		req.Address,
		req.IsActive,
	); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Update(r.Context(), location); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Lokasi berhasil diupdate", toLocationOutput(location))
}

func (h *LocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	_, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, "Lokasi tidak ditemukan")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Lokasi berhasil dihapus", nil)
}

func toLocationOutput(l *entity.AllowedLocation) *LocationOutput {
	return &LocationOutput{
		ID:           l.ID,
		Name:         l.Name,
		Latitude:     l.Latitude,
		Longitude:    l.Longitude,
		RadiusMeters: l.RadiusMeters,
		Address:      l.Address,
		IsActive:     l.IsActive,
		CreatedAt:    l.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    l.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
