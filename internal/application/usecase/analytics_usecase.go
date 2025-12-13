package usecase

import (
	"context"
	"time"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

// AnalyticsUseCase adalah interface untuk analytics use case
type AnalyticsUseCase interface {
	GetSummary(ctx context.Context, startDate, endDate string) (*entity.AttendanceSummary, error)
	GetDailySummary(ctx context.Context, date string) (*entity.DailySummary, error)
	GetMonthlySummary(ctx context.Context, month string) (*entity.MonthlySummary, error)
	GetUserSummary(ctx context.Context, userID, startDate, endDate string) (*entity.UserSummary, error)
	GetStatusBreakdown(ctx context.Context, startDate, endDate string) ([]entity.StatusBreakdown, error)
}

type analyticsUseCase struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsUseCase(repo repository.AnalyticsRepository) AnalyticsUseCase {
	return &analyticsUseCase{
		repo: repo,
	}
}

func (uc *analyticsUseCase) GetSummary(ctx context.Context, startDate, endDate string) (*entity.AttendanceSummary, error) {
	filter := entity.AnalyticsFilter{}

	if startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		filter.StartDate = parsed
	}

	if endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		// Include the entire end date
		filter.EndDate = parsed.Add(24 * time.Hour)
	}

	return uc.repo.GetSummary(ctx, filter)
}

func (uc *analyticsUseCase) GetDailySummary(ctx context.Context, date string) (*entity.DailySummary, error) {
	return uc.repo.GetDailySummary(ctx, date)
}

func (uc *analyticsUseCase) GetMonthlySummary(ctx context.Context, month string) (*entity.MonthlySummary, error) {
	return uc.repo.GetMonthlySummary(ctx, month)
}

func (uc *analyticsUseCase) GetUserSummary(ctx context.Context, userID, startDate, endDate string) (*entity.UserSummary, error) {
	filter := entity.AnalyticsFilter{}

	if startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		filter.StartDate = parsed
	}

	if endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		filter.EndDate = parsed.Add(24 * time.Hour)
	}

	return uc.repo.GetUserSummary(ctx, userID, filter)
}

func (uc *analyticsUseCase) GetStatusBreakdown(ctx context.Context, startDate, endDate string) ([]entity.StatusBreakdown, error) {
	filter := entity.AnalyticsFilter{}

	if startDate != "" {
		parsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		filter.StartDate = parsed
	}

	if endDate != "" {
		parsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		filter.EndDate = parsed.Add(24 * time.Hour)
	}

	return uc.repo.GetStatusBreakdown(ctx, filter)
}
