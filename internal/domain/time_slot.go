package domain

import "time"

type TimeSlot struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Time      time.Time `json:"time"`
	IsBooked  bool      `json:"is_booked"`
	CreatedAt time.Time `json:"created_at"`
}
