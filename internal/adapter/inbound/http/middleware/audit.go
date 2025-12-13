/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

// AuditContextKey is the context key for audit information
type AuditContextKey string

const (
	ContextKeyRequestID   AuditContextKey = "request_id"
	ContextKeyAuditLog    AuditContextKey = "audit_log"
	ContextKeyRequestBody AuditContextKey = "request_body"
)

// AuditMiddleware creates middleware that logs all API requests for audit trail
type AuditMiddleware struct {
	auditRepo repository.AuditLogRepository
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(auditRepo repository.AuditLogRepository) *AuditMiddleware {
	return &AuditMiddleware{
		auditRepo: auditRepo,
	}
}

// auditResponseWriter wraps http.ResponseWriter to capture status code and response
type auditResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseBody *bytes.Buffer
}

func newAuditResponseWriter(w http.ResponseWriter) *auditResponseWriter {
	return &auditResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		responseBody:   &bytes.Buffer{},
	}
}

func (w *auditResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *auditResponseWriter) Write(data []byte) (int, error) {
	w.responseBody.Write(data)
	return w.ResponseWriter.Write(data)
}

// Audit returns the audit middleware handler
func (m *AuditMiddleware) Audit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip audit for certain paths
		if shouldSkipAudit(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()
		requestID := generateRequestID()

		// Read and restore request body
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Add request ID and body to context
		ctx := context.WithValue(r.Context(), ContextKeyRequestID, requestID)
		ctx = context.WithValue(ctx, ContextKeyRequestBody, string(requestBody))
		r = r.WithContext(ctx)

		// Wrap response writer
		wrappedWriter := newAuditResponseWriter(w)

		// Add request ID to response header
		wrappedWriter.Header().Set("X-Request-ID", requestID)

		// Call next handler
		next.ServeHTTP(wrappedWriter, r)

		// Calculate duration
		duration := time.Since(startTime).Milliseconds()

		// Create audit log asynchronously
		go m.createAuditLog(r, requestBody, wrappedWriter, requestID, duration)
	})
}

func (m *AuditMiddleware) createAuditLog(
	r *http.Request,
	requestBody []byte,
	w *auditResponseWriter,
	requestID string,
	duration int64,
) {
	// Determine entity type and action from request
	entityType, entityID := parseEntityFromPath(r.URL.Path)
	action := determineAction(r.Method)

	if entityType == "" || action == "" {
		return
	}

	// Get user info from context (set by auth middleware)
	userID, userName, userRole := getUserFromContext(r.Context())

	// Create audit log
	auditLog := entity.NewAuditLog(
		entityType,
		entityID,
		action,
		userID,
		userName,
		userRole,
		getClientIP(r),
	)

	// Set request info
	auditLog.SetRequestInfo(
		requestID,
		r.UserAgent(),
		w.statusCode,
		duration,
	)

	// Set changes (sanitize sensitive data)
	sanitizedBody := sanitizeRequestBody(requestBody)
	auditLog.SetChanges("", sanitizedBody, "")

	// Check for errors in response
	if w.statusCode >= 400 {
		auditLog.SetError(extractErrorFromResponse(w.responseBody.Bytes()))
	}

	// Save to database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m.auditRepo.Create(ctx, auditLog)
}

// Helper functions

func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/health",
		"/metrics",
		"/favicon.ico",
	}

	for _, p := range skipPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	// Only audit write operations for GET
	return false
}

func generateRequestID() string {
	return uuid.New().String()
}

func parseEntityFromPath(path string) (entityType, entityID string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Expected format: /api/{entity}/{id}
	if len(parts) < 2 {
		return "", ""
	}

	// Skip "api" prefix
	startIdx := 0
	if parts[0] == "api" {
		startIdx = 1
	}

	if len(parts) > startIdx {
		entityType = parts[startIdx]
	}

	if len(parts) > startIdx+1 {
		entityID = parts[startIdx+1]
	}

	return entityType, entityID
}

func determineAction(method string) entity.AuditAction {
	switch method {
	case http.MethodPost:
		return entity.AuditActionCreate
	case http.MethodPut, http.MethodPatch:
		return entity.AuditActionUpdate
	case http.MethodDelete:
		return entity.AuditActionDelete
	default:
		return ""
	}
}

func getUserFromContext(ctx context.Context) (userID, userName, userRole string) {
	// These should be set by auth middleware
	if id := ctx.Value("user_id"); id != nil {
		userID = id.(string)
	}
	if name := ctx.Value("user_name"); name != nil {
		userName = name.(string)
	}
	if role := ctx.Value("user_role"); role != nil {
		userRole = role.(string)
	}
	return
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colonIdx := strings.LastIndex(ip, ":"); colonIdx != -1 {
		ip = ip[:colonIdx]
	}
	return ip
}

func sanitizeRequestBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	// Parse JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "[non-json body]"
	}

	// Remove sensitive fields
	sensitiveFields := []string{"password", "token", "secret", "api_key", "credit_card"}
	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "[REDACTED]"
		}
	}

	// Re-serialize
	sanitized, _ := json.Marshal(data)
	return string(sanitized)
}

func extractErrorFromResponse(body []byte) string {
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return ""
	}

	if errMsg, ok := response["error"].(string); ok {
		return errMsg
	}
	if errMsg, ok := response["message"].(string); ok {
		return errMsg
	}
	return ""
}
