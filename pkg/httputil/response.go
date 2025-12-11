package httputil

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(w http.ResponseWriter, status int, message string, data interface{}, meta *Meta) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{
		Success: false,
		Message: message,
	})
}
