package entity

import "time"

type AttendanceSummary struct {
	TotalRecords    int     `json:"total_records"`
	TotalHadir      int     `json:"total_hadir"`
	TotalTerlambat  int     `json:"total_terlambat"`
	TotalIzin       int     `json:"total_izin"`
	TotalSakit      int     `json:"total_sakit"`
	TotalAlpha      int     `json:"total_alpha"`
	PercentageHadir float64 `json:"percentage_hadir"`
}

type StatusBreakdown struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type DailySummary struct {
	Date    time.Time          `json:"date"`
	Summary AttendanceSummary  `json:"summary"`
	Details []StatusBreakdown  `json:"details"`
}

type MonthlySummary struct {
	Month       string             `json:"month"` // Format: YYYY-MM
	Summary     AttendanceSummary  `json:"summary"`
	DailyStats  []DailyStats       `json:"daily_stats,omitempty"`
}

type DailyStats struct {
	Date  string `json:"date"` // Format: YYYY-MM-DD
	Count int    `json:"count"`
}

type UserSummary struct {
	UserID       string            `json:"user_id"`
	UserName     string            `json:"user_name"`
	Period       string            `json:"period"` // e.g., "2024-01" or "2024-01-01 to 2024-01-31"
	Summary      AttendanceSummary `json:"summary"`
	StatusDetail []StatusBreakdown `json:"status_detail"`
}

// AnalyticsFilter for filtering analytics queries
type AnalyticsFilter struct {
	UserID    string
	StartDate time.Time
	EndDate   time.Time
}

func (s *AttendanceSummary) CalculatePercentage() {
	if s.TotalRecords > 0 {
		// Hadir includes both "hadir" and "terlambat"
		totalPresent := s.TotalHadir + s.TotalTerlambat
		s.PercentageHadir = float64(totalPresent) / float64(s.TotalRecords) * 100
	}
}
