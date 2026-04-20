package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bookify/internal/domain"
	"github.com/bookify/pkg/validator"
)

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details []validator.FieldError `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "An internal error occurred"

	switch {
	case errors.As(err, new(validator.ValidationErrors)):
		var validationErrs validator.ValidationErrors
		_ = errors.As(err, &validationErrs)
		writeJSON(w, http.StatusBadRequest, errorEnvelope{
			Error: apiError{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid input data",
				Details: []validator.FieldError(validationErrs),
			},
		})
		return
	case errors.Is(err, domain.ErrValidation):
		status = http.StatusBadRequest
		code = "VALIDATION_ERROR"
		message = "Invalid input data"
	case errors.Is(err, domain.ErrAlreadyExists):
		status = http.StatusBadRequest
		code = "ALREADY_EXISTS"
		message = err.Error()
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
		code = "UNAUTHORIZED"
		message = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		code = "FORBIDDEN"
		message = err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		code = "NOT_FOUND"
		message = err.Error()
	case errors.Is(err, domain.ErrTimeSlotTaken), errors.Is(err, domain.ErrHasFutureBookings):
		status = http.StatusConflict
		code = "CONFLICT"
		message = err.Error()
	case errors.Is(err, domain.ErrPastTime), errors.Is(err, domain.ErrServiceInactive), errors.Is(err, domain.ErrInvalidTimeSlot):
		status = http.StatusBadRequest
		code = "BUSINESS_RULE_ERROR"
		message = err.Error()
	default:
		if err != nil {
			status = http.StatusBadRequest
			code = "BAD_REQUEST"
			message = err.Error()
		}
	}

	writeJSON(w, status, errorEnvelope{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}
