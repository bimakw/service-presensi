package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/okinn/service-presensi/internal/application/usecase"
	"github.com/okinn/service-presensi/internal/domain/repository"
	"github.com/okinn/service-presensi/internal/domain/valueobject"
	"github.com/okinn/service-presensi/pkg/validator"
)

type PresensiHandler struct {
	useCase usecase.PresensiUseCase
}

func NewPresensiHandler(uc usecase.PresensiUseCase) *PresensiHandler {
	return &PresensiHandler{useCase: uc}
}

type CreatePresensiRequest struct {
	UserID     string  `json:"user_id" validate:"required"`
	Nama       string  `json:"nama" validate:"required,min=2,max=100"`
	Status     string  `json:"status" validate:"required,status_presensi"`
	Keterangan string  `json:"keterangan" validate:"max=500"`
	Latitude   float64 `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude  float64 `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	Alamat     string  `json:"alamat" validate:"max=255"`
}

type UpdatePresensiRequest struct {
	Status     string `json:"status" validate:"omitempty,status_presensi"`
	Keterangan string `json:"keterangan" validate:"max=500"`
}

func (h *PresensiHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePresensiRequest
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

	input := usecase.CreatePresensiInput{
		UserID:     req.UserID,
		Nama:       req.Nama,
		Status:     req.Status,
		Keterangan: req.Keterangan,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Alamat:     req.Alamat,
	}

	output, err := h.useCase.Create(r.Context(), input)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusCreated, "Presensi berhasil dibuat", output)
}

func (h *PresensiHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	output, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil", output)
}

func (h *PresensiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filter := repository.PresensiFilter{
		UserID: query.Get("user_id"),
		Status: valueobject.StatusPresensi(query.Get("status")),
	}

	if startDate := query.Get("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = t
		}
	}

	if endDate := query.Get("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = t.Add(24*time.Hour - time.Second)
		}
	}

	outputs, total, err := h.useCase.GetAll(r.Context(), filter, page, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	SuccessWithMeta(w, http.StatusOK, "Berhasil", outputs, &Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *PresensiHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var req UpdatePresensiRequest
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

	input := usecase.UpdatePresensiInput{
		Status:     req.Status,
		Keterangan: req.Keterangan,
	}

	output, err := h.useCase.Update(r.Context(), id, input)
	if err != nil {
		if err == usecase.ErrPresensiNotFound {
			Error(w, http.StatusNotFound, err.Error())
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Presensi berhasil diupdate", output)
}

func (h *PresensiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	err := h.useCase.Delete(r.Context(), id)
	if err != nil {
		if err == usecase.ErrPresensiNotFound {
			Error(w, http.StatusNotFound, err.Error())
			return
		}
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Success(w, http.StatusOK, "Presensi berhasil dihapus", nil)
}

func (h *PresensiHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	err := h.useCase.CheckIn(r.Context(), id)
	if err != nil {
		if err == usecase.ErrPresensiNotFound {
			Error(w, http.StatusNotFound, err.Error())
			return
		}
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusOK, "Check-in berhasil", nil)
}

func (h *PresensiHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	err := h.useCase.CheckOut(r.Context(), id)
	if err != nil {
		if err == usecase.ErrPresensiNotFound {
			Error(w, http.StatusNotFound, err.Error())
			return
		}
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Success(w, http.StatusOK, "Check-out berhasil", nil)
}
