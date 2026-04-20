package domain

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrValidation         = errors.New("validation error")
	ErrTimeSlotTaken      = errors.New("time slot is already booked")
	ErrInvalidTimeSlot    = errors.New("invalid time slot")
	ErrPastTime           = errors.New("cannot book appointment in the past")
	ErrServiceInactive    = errors.New("service is not active")
	ErrHasFutureBookings  = errors.New("cannot delete service with future bookings")
)
