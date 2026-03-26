package repository

import (
	"database/sql"

	"github.com/bookify/internal/domain"
)

type TimeSlotRepository struct {
	db *sql.DB
}

func NewTimeSlotRepository(db *sql.DB) *TimeSlotRepository {
	return &TimeSlotRepository{db: db}
}

func (r *TimeSlotRepository) GetBySpecialistID(specialistID int) ([]domain.TimeSlot, error) {
	query := `SELECT id, specialist_id, time, is_booked FROM time_slots WHERE specialist_id = $1 ORDER BY time`

	rows, err := r.db.Query(query, specialistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []domain.TimeSlot
	for rows.Next() {
		var slot domain.TimeSlot
		err := rows.Scan(&slot.ID, &slot.SpecialistID, &slot.Time, &slot.IsBooked)
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

func (r *TimeSlotRepository) GetByID(id int) (*domain.TimeSlot, error) {
	query := `SELECT id, specialist_id, time, is_booked FROM time_slots WHERE id = $1`

	var slot domain.TimeSlot
	err := r.db.QueryRow(query, id).Scan(&slot.ID, &slot.SpecialistID, &slot.Time, &slot.IsBooked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &slot, nil
}

func (r *TimeSlotRepository) MarkAsBooked(id int) error {
	query := `UPDATE time_slots SET is_booked = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
