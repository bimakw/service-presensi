/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package http

import (
	"net/http"
	"strconv"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

// AuditHandler handles HTTP requests for audit logs
type AuditHandler struct {
	auditRepo repository.AuditLogRepository
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(auditRepo repository.AuditLogRepository) *AuditHandler {
	return &AuditHandler{
		auditRepo: auditRepo,
	}
}

// GetAll returns all audit logs with filtering and pagination
func (h *AuditHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Build filter
	filter := repository.AuditLogFilter{
		EntityType: r.URL.Query().Get("entity_type"),
		EntityID:   r.URL.Query().Get("entity_id"),
		UserID:     r.URL.Query().Get("user_id"),
		IPAddress:  r.URL.Query().Get("ip_address"),
	}

	if action := r.URL.Query().Get("action"); action != "" {
		filter.Action = entity.AuditAction(action)
	}

	// Get audit logs
	auditLogs, total, err := h.auditRepo.GetAll(ctx, filter, page, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Gagal mengambil audit logs")
		return
	}

	// Build pagination response
	response := map[string]interface{}{
		"data":  auditLogs,
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": (total + int64(limit) - 1) / int64(limit),
	}

	Success(w, http.StatusOK, "Berhasil mengambil audit logs", response)
}

// GetByID returns a single audit log by ID
func (h *AuditHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	auditLog, err := h.auditRepo.GetByID(ctx, id)
	if err != nil {
		Error(w, http.StatusNotFound, "Audit log tidak ditemukan")
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil audit log", auditLog)
}

// GetByEntity returns audit logs for a specific entity
func (h *AuditHandler) GetByEntity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	entityType := r.URL.Query().Get("type")
	entityID := r.URL.Query().Get("id")

	if entityType == "" || entityID == "" {
		Error(w, http.StatusBadRequest, "Parameter type dan id wajib diisi")
		return
	}

	auditLogs, err := h.auditRepo.GetByEntityID(ctx, entityType, entityID)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Gagal mengambil audit logs")
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil audit logs", auditLogs)
}

// GetByUser returns audit logs for a specific user
func (h *AuditHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.PathValue("user_id")

	if userID == "" {
		Error(w, http.StatusBadRequest, "User ID tidak valid")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	auditLogs, total, err := h.auditRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Gagal mengambil audit logs")
		return
	}

	response := map[string]interface{}{
		"data":  auditLogs,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	Success(w, http.StatusOK, "Berhasil mengambil audit logs", response)
}
