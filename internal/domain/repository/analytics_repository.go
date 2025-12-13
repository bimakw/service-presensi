package repository

import (
	"context"

	"github.com/okinn/service-presensi/internal/domain/entity"
)

// AnalyticsRepository adalah port untuk analytics data
type AnalyticsRepository interface {
	// GetSummary returns overall attendance summary with optional date filter
	GetSummary(ctx context.Context, filter entity.AnalyticsFilter) (*entity.AttendanceSummary, error)

	// GetDailySummary returns attendance summary for a specific date
	GetDailySummary(ctx context.Context, date string) (*entity.DailySummary, error)

	// GetMonthlySummary returns attendance summary for a month
	GetMonthlySummary(ctx context.Context, month string) (*entity.MonthlySummary, error)

	// GetUserSummary returns attendance summary for a specific user
	GetUserSummary(ctx context.Context, userID string, filter entity.AnalyticsFilter) (*entity.UserSummary, error)

	// GetStatusBreakdown returns count per status
	GetStatusBreakdown(ctx context.Context, filter entity.AnalyticsFilter) ([]entity.StatusBreakdown, error)
}
