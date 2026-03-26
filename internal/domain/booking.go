package domain

import "time"

type Booking struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TimeSlotID int       `json:"time_slot_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type BookingResponse struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TimeSlotID int       `json:"time_slot_id"`
	CreatedAt  time.Time `json:"created_at"`
	Specialist string    `json:"specialist"`
	SlotTime   time.Time `json:"slot_time"`
}
