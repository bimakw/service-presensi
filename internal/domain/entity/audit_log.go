/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package entity

import (
	"time"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
	AuditActionLogin  AuditAction = "login"
	AuditActionLogout AuditAction = "logout"
)

// AuditLog represents an audit trail entry for tracking all data changes
type AuditLog struct {
	ID          string      `json:"id"`
	EntityType  string      `json:"entity_type"`  // "presensi", "user", etc.
	EntityID    string      `json:"entity_id"`    // ID of the affected entity
	Action      AuditAction `json:"action"`       // create, update, delete
	OldValue    string      `json:"old_value"`    // JSON string of old data (for updates)
	NewValue    string      `json:"new_value"`    // JSON string of new data
	Changes     string      `json:"changes"`      // JSON diff of changes
	UserID      string      `json:"user_id"`      // User who performed the action
	UserName    string      `json:"user_name"`    // Username for display
	UserRole    string      `json:"user_role"`    // Role at time of action
	IPAddress   string      `json:"ip_address"`   // Client IP address
	UserAgent   string      `json:"user_agent"`   // Browser/client info
	RequestID   string      `json:"request_id"`   // For request tracing
	StatusCode  int         `json:"status_code"`  // HTTP status code
	ErrorMsg    string      `json:"error_msg"`    // Error message if failed
	Duration    int64       `json:"duration_ms"`  // Request duration in ms
	CreatedAt   time.Time   `json:"created_at"`
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(
	entityType string,
	entityID string,
	action AuditAction,
	userID string,
	userName string,
	userRole string,
	ipAddress string,
) *AuditLog {
	return &AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		UserID:     userID,
		UserName:   userName,
		UserRole:   userRole,
		IPAddress:  ipAddress,
		CreatedAt:  time.Now(),
	}
}

// SetChanges sets the old value, new value, and computed changes
func (a *AuditLog) SetChanges(oldValue, newValue, changes string) {
	a.OldValue = oldValue
	a.NewValue = newValue
	a.Changes = changes
}

// SetRequestInfo sets request-related information
func (a *AuditLog) SetRequestInfo(requestID, userAgent string, statusCode int, duration int64) {
	a.RequestID = requestID
	a.UserAgent = userAgent
	a.StatusCode = statusCode
	a.Duration = duration
}

// SetError sets error information for failed operations
func (a *AuditLog) SetError(errMsg string) {
	a.ErrorMsg = errMsg
}

// IsSuccessful returns true if the operation was successful
func (a *AuditLog) IsSuccessful() bool {
	return a.ErrorMsg == "" && a.StatusCode >= 200 && a.StatusCode < 300
}
