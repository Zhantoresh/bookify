package domain

import "time"

type TimeSlot struct {
	ID           int       `json:"id"`
	SpecialistID int       `json:"specialist_id"`
	Time         time.Time `json:"time"`
	IsBooked     bool      `json:"is_booked"`
}
