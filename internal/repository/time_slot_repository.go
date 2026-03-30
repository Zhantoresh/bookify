package repository

import (
	"database/sql"
	"time"

	"github.com/bookify/internal/domain"
)

type TimeSlotRepository struct {
	db *sql.DB
}

func NewTimeSlotRepository(db *sql.DB) *TimeSlotRepository {
	return &TimeSlotRepository{db: db}
}

// GetByUserID returns all time slots for a specialist (user with role specialist)
func (r *TimeSlotRepository) GetByUserID(userID int) ([]domain.TimeSlot, error) {
	query := `SELECT id, user_id, time, is_booked, created_at FROM time_slots WHERE user_id = $1 ORDER BY time`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []domain.TimeSlot
	for rows.Next() {
		var slot domain.TimeSlot
		err := rows.Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsBooked, &slot.CreatedAt)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

// GetBySpecialistID for backward compatibility
func (r *TimeSlotRepository) GetBySpecialistID(specialistID int) ([]domain.TimeSlot, error) {
	return r.GetByUserID(specialistID)
}

func (r *TimeSlotRepository) GetByID(id int) (*domain.TimeSlot, error) {
	query := `SELECT id, user_id, time, is_booked, created_at FROM time_slots WHERE id = $1`

	var slot domain.TimeSlot
	err := r.db.QueryRow(query, id).Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsBooked, &slot.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &slot, nil
}

// CreateTimeSlot creates a new time slot for a specialist
func (r *TimeSlotRepository) CreateTimeSlot(userID int, slotTime time.Time) (*domain.TimeSlot, error) {
	slot := &domain.TimeSlot{}
	query := `INSERT INTO time_slots (user_id, time, is_booked) VALUES ($1, $2, false) RETURNING id, user_id, time, is_booked, created_at`
	err := r.db.QueryRow(query, userID, slotTime).Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsBooked, &slot.CreatedAt)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

// UpdateTimeSlot updates a time slot
func (r *TimeSlotRepository) UpdateTimeSlot(id int, slotTime time.Time) error {
	query := `UPDATE time_slots SET time = $1 WHERE id = $2`
	_, err := r.db.Exec(query, slotTime, id)
	return err
}

// DeleteTimeSlot deletes a time slot
func (r *TimeSlotRepository) DeleteTimeSlot(id int) error {
	query := `DELETE FROM time_slots WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TimeSlotRepository) MarkAsBooked(id int) error {
	query := `UPDATE time_slots SET is_booked = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TimeSlotRepository) MarkAsUnbooked(id int) error {
	query := `UPDATE time_slots SET is_booked = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetAvailableSlots returns all available (not booked) time slots
func (r *TimeSlotRepository) GetAvailableSlots() ([]domain.TimeSlot, error) {
	query := `SELECT id, user_id, time, is_booked, created_at FROM time_slots WHERE is_booked = false ORDER BY time`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []domain.TimeSlot
	for rows.Next() {
		var slot domain.TimeSlot
		err := rows.Scan(&slot.ID, &slot.UserID, &slot.Time, &slot.IsBooked, &slot.CreatedAt)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
