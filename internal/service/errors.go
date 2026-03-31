package service

import "errors"

var (
	ErrTimeSlotNotFound        = errors.New("time slot not found")
	ErrBookingAlreadyBooked    = errors.New("this slot is already booked")
	ErrSpecialistNotFound      = errors.New("specialist not found")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrForbidden               = errors.New("forbidden")
	ErrBookingAlreadyCancelled = errors.New("booking is already cancelled")
)
