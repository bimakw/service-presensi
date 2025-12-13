package entity

import "time"

// AttendanceSummary represents aggregated attendance statistics
type AttendanceSummary struct {
	TotalRecords    int     `json:"total_records"`
	TotalHadir      int     `json:"total_hadir"`
	TotalTerlambat  int     `json:"total_terlambat"`
	TotalIzin       int     `json:"total_izin"`
	TotalSakit      int     `json:"total_sakit"`
	TotalAlpha      int     `json:"total_alpha"`
	PercentageHadir float64 `json:"percentage_hadir"`
}

// StatusBreakdown represents count per status
type StatusBreakdown struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// DailySummary represents attendance summary for a specific date
type DailySummary struct {
	Date    time.Time          `json:"date"`
	Summary AttendanceSummary  `json:"summary"`
	Details []StatusBreakdown  `json:"details"`
}

// MonthlySummary represents attendance summary for a month
type MonthlySummary struct {
	Month       string             `json:"month"` // Format: YYYY-MM
	Summary     AttendanceSummary  `json:"summary"`
	DailyStats  []DailyStats       `json:"daily_stats,omitempty"`
}

// DailyStats represents daily count within a month
type DailyStats struct {
	Date  string `json:"date"` // Format: YYYY-MM-DD
	Count int    `json:"count"`
}

// UserSummary represents attendance summary for a specific user
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

// CalculatePercentage calculates and sets the attendance percentage
func (s *AttendanceSummary) CalculatePercentage() {
	if s.TotalRecords > 0 {
		// Hadir includes both "hadir" and "terlambat"
		totalPresent := s.TotalHadir + s.TotalTerlambat
		s.PercentageHadir = float64(totalPresent) / float64(s.TotalRecords) * 100
	}
}
