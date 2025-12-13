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
	"time"

	"github.com/okinn/service-presensi/internal/domain/entity"
)

// AuditLogFilter untuk filtering audit logs
type AuditLogFilter struct {
	EntityType string
	EntityID   string
	Action     entity.AuditAction
	UserID     string
	StartDate  time.Time
	EndDate    time.Time
	IPAddress  string
}

// AuditLogRepository adalah port untuk akses data audit log
type AuditLogRepository interface {
	// Create menyimpan audit log baru
	Create(ctx context.Context, auditLog *entity.AuditLog) error

	// GetByID mengambil audit log berdasarkan ID
	GetByID(ctx context.Context, id string) (*entity.AuditLog, error)

	// GetAll mengambil semua audit logs dengan filter dan pagination
	GetAll(ctx context.Context, filter AuditLogFilter, page, limit int) ([]entity.AuditLog, int64, error)

	// GetByEntityID mengambil audit logs untuk entity tertentu
	GetByEntityID(ctx context.Context, entityType, entityID string) ([]entity.AuditLog, error)

	// GetByUserID mengambil audit logs untuk user tertentu
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]entity.AuditLog, int64, error)

	// DeleteOldLogs menghapus audit logs yang lebih lama dari retention period
	DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error)
}
