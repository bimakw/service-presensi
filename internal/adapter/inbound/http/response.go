package http

import (
	"net/http"

	"github.com/okinn/service-presensi/pkg/httputil"
	"github.com/okinn/service-presensi/pkg/validator"
)

// Re-export types from httputil for backward compatibility
type Response = httputil.Response
type Meta = httputil.Meta

func JSON(w http.ResponseWriter, status int, data interface{}) {
	httputil.JSON(w, status, data)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	httputil.Success(w, status, message, data)
}

func SuccessWithMeta(w http.ResponseWriter, status int, message string, data interface{}, meta *Meta) {
	httputil.SuccessWithMeta(w, status, message, data, meta)
}

func Error(w http.ResponseWriter, status int, message string) {
	httputil.Error(w, status, message)
}

func ValidationError(w http.ResponseWriter, errors validator.ValidationErrors) {
	httputil.JSON(w, http.StatusBadRequest, Response{
		Success: false,
		Message: "Validasi gagal",
		Data:    errors,
	})
}
