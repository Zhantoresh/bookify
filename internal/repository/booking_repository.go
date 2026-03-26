package repository

import (
	"database/sql"
	"time"

	"github.com/bookify/internal/domain"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(userID, timeSlotID int) (*domain.Booking, error) {
	query := `INSERT INTO bookings (user_id, time_slot_id, status, created_at) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, user_id, time_slot_id, status, created_at`

	var booking domain.Booking
	err := r.db.QueryRow(query, userID, timeSlotID, "BOOKED", time.Now()).
		Scan(&booking.ID, &booking.UserID, &booking.TimeSlotID, &booking.Status, &booking.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) GetByID(id int) (*domain.Booking, error) {
	query := `SELECT id, user_id, time_slot_id, status, created_at FROM bookings WHERE id = $1`

	var booking domain.Booking
	err := r.db.QueryRow(query, id).Scan(&booking.ID, &booking.UserID, &booking.TimeSlotID, &booking.Status, &booking.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) GetByUserID(userID int) ([]domain.Booking, error) {
	query := `SELECT id, user_id, time_slot_id, status, created_at 
	          FROM bookings 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.TimeSlotID, &booking.Status, &booking.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) CancelByID(id int) error {
	query := `UPDATE bookings SET status = 'CANCELLED' WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
