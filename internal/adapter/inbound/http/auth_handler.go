package http

import (
	"encoding/json"
	"net/http"

	"github.com/okinn/service-presensi/internal/adapter/inbound/http/middleware"
	"github.com/okinn/service-presensi/internal/application/usecase"
	"github.com/okinn/service-presensi/pkg/validator"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: uc}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nama     string `json:"nama" validate:"required,min=2,max=100"`
	Role     string `json:"role" validate:"omitempty,role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
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

	input := usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Nama:     req.Nama,
		Role:     req.Role,
	}

	output, err := h.useCase.Register(r.Context(), input)
	if err != nil {
		switch err {
		case usecase.ErrEmailAlreadyExists:
			Error(w, http.StatusConflict, err.Error())
		default:
			Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	Success(w, http.StatusCreated, "Registrasi berhasil", output)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
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

	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.useCase.Login(r.Context(), input)
	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			Error(w, http.StatusUnauthorized, err.Error())
		case usecase.ErrUserNotActive:
			Error(w, http.StatusForbidden, err.Error())
		default:
			Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Success(w, http.StatusOK, "Login berhasil", output)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		Error(w, http.StatusUnauthorized, "User ID tidak ditemukan")
		return
	}

	output, err := h.useCase.GetProfile(r.Context(), userID)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Success(w, http.StatusOK, "Berhasil", output)
}
