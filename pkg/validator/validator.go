package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function untuk menggunakan json tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	validate.RegisterValidation("status_presensi", validateStatusPresensi)
	validate.RegisterValidation("role", validateRole)
}

// Custom validation untuk status presensi
func validateStatusPresensi(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"hadir", "terlambat", "izin", "sakit", "alpha"}
	for _, s := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

// Custom validation untuk role
func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	if role == "" {
		return true // empty role is allowed (will default to employee)
	}
	validRoles := []string{"admin", "employee"}
	for _, r := range validRoles {
		if role == r {
			return true
		}
	}
	return false
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Validate validates a struct and returns formatted errors
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = "field ini wajib diisi"
		case "email":
			message = "format email tidak valid"
		case "min":
			message = fmt.Sprintf("minimal %s karakter", err.Param())
		case "max":
			message = fmt.Sprintf("maksimal %s karakter", err.Param())
		case "status_presensi":
			message = "status harus salah satu dari: hadir, terlambat, izin, sakit, alpha"
		case "role":
			message = "role harus salah satu dari: admin, employee"
		case "latitude":
			message = "latitude harus antara -90 dan 90"
		case "longitude":
			message = "longitude harus antara -180 dan 180"
		case "gte":
			message = fmt.Sprintf("nilai minimal %s", err.Param())
		case "lte":
			message = fmt.Sprintf("nilai maksimal %s", err.Param())
		default:
			message = fmt.Sprintf("validasi %s gagal", err.Tag())
		}

		errors = append(errors, ValidationError{
			Field:   err.Field(),
			Message: message,
		})
	}

	return errors
}
