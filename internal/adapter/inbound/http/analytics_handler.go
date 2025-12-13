package http

import (
	"net/http"

	"github.com/okinn/service-presensi/internal/application/usecase"
)

type AnalyticsHandler struct {
	useCase usecase.AnalyticsUseCase
}

func NewAnalyticsHandler(uc usecase.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{useCase: uc}
}

// GetSummary returns overall attendance summary
// GET /api/analytics/summary?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	summary, err := h.useCase.GetSummary(r.Context(), startDate, endDate)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil summary", summary)
}

// GetDailySummary returns attendance summary for a specific date
// GET /api/analytics/daily?date=YYYY-MM-DD
func (h *AnalyticsHandler) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		Error(w, http.StatusBadRequest, "Parameter 'date' diperlukan (format: YYYY-MM-DD)")
		return
	}

	summary, err := h.useCase.GetDailySummary(r.Context(), date)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil summary harian", summary)
}

// GetMonthlySummary returns attendance summary for a month
// GET /api/analytics/monthly?month=YYYY-MM
func (h *AnalyticsHandler) GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	if month == "" {
		Error(w, http.StatusBadRequest, "Parameter 'month' diperlukan (format: YYYY-MM)")
		return
	}

	summary, err := h.useCase.GetMonthlySummary(r.Context(), month)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil summary bulanan", summary)
}

// GetUserSummary returns attendance summary for a specific user
// GET /api/analytics/user/{user_id}?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *AnalyticsHandler) GetUserSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		Error(w, http.StatusBadRequest, "User ID diperlukan")
		return
	}

	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	summary, err := h.useCase.GetUserSummary(r.Context(), userID, startDate, endDate)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil summary user", summary)
}

// GetStatusBreakdown returns count per status
// GET /api/analytics/status-breakdown?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *AnalyticsHandler) GetStatusBreakdown(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	breakdown, err := h.useCase.GetStatusBreakdown(r.Context(), startDate, endDate)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil mengambil status breakdown", breakdown)
}
