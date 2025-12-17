package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/okinn/service-presensi/internal/adapter/inbound/http/middleware"
	"github.com/okinn/service-presensi/internal/application/usecase"
	"github.com/okinn/service-presensi/pkg/validator"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	output, err := h.useCase.GetAllUsers(r.Context(), page, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (output.Total + int64(limit) - 1) / int64(limit)
	SuccessWithMeta(w, http.StatusOK, "Berhasil", output.Users, &Meta{
		Page:       page,
		Limit:      limit,
		Total:      output.Total,
		TotalPages: totalPages,
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "User ID tidak boleh kosong")
		return
	}

	output, err := h.useCase.GetUserByID(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil", output)
}

func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		Error(w, http.StatusBadRequest, "User ID tidak boleh kosong")
		return
	}

	// Prevent admin from deactivating themselves
	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID == id {
		Error(w, http.StatusBadRequest, "Tidak dapat mengubah status diri sendiri")
		return
	}

	var req UpdateUserStatusRequest
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

	err := h.useCase.UpdateUserStatus(r.Context(), id, req.IsActive)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			Error(w, http.StatusNotFound, err.Error())
		default:
			Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	message := "User berhasil diaktifkan"
	if !req.IsActive {
		message = "User berhasil dinonaktifkan"
	}

	Success(w, http.StatusOK, message, nil)
}
