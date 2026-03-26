package domain

import "time"

type Booking struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TimeSlotID int       `json:"time_slot_id"`
	Status     string    `json:"status"` // BOOKED or CANCELLED
	CreatedAt  time.Time `json:"created_at"`
}

type BookingResponse struct {
	ID         int       `json:"id"`
	Specialist string    `json:"specialist"`
	Time       time.Time `json:"time"`
	Status     string    `json:"status"`
}

type CreateBookingResponse struct {
	ID         int       `json:"id"`
	Specialist string    `json:"specialist"`
	Time       time.Time `json:"time"`
	Status     string    `json:"status"`
}
