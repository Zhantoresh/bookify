package service

import (
	"errors"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type BookingService struct {
	bookingRepo    *repository.BookingRepository
	timeSlotRepo   *repository.TimeSlotRepository
	specialistRepo *repository.SpecialistRepository
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	timeSlotRepo *repository.TimeSlotRepository,
	specialistRepo *repository.SpecialistRepository,
) *BookingService {
	return &BookingService{
		bookingRepo:    bookingRepo,
		timeSlotRepo:   timeSlotRepo,
		specialistRepo: specialistRepo,
	}
}

func (s *BookingService) CreateBooking(userID, timeSlotID int) (*domain.Booking, error) {
	// Check if time slot exists
	timeSlot, err := s.timeSlotRepo.GetByID(timeSlotID)
	if err != nil {
		return nil, err
	}

	if timeSlot == nil {
		return nil, errors.New("time slot not found")
	}

	// Check if slot is already booked
	if timeSlot.IsBooked {
		return nil, errors.New("this slot is already booked")
	}

	booking, err := s.bookingRepo.Create(userID, timeSlotID)
	if err != nil {
		return nil, err
	}

	err = s.timeSlotRepo.MarkAsBooked(timeSlotID)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) GetUserBookings(userID int) ([]domain.BookingResponse, error) {
	bookings, err := s.bookingRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []domain.BookingResponse

	for _, booking := range bookings {
		timeSlot, err := s.timeSlotRepo.GetByID(booking.TimeSlotID)
		if err != nil {
			return nil, err
		}

		if timeSlot == nil {
			continue
		}

		specialist, err := s.specialistRepo.GetByID(timeSlot.SpecialistID)
		if err != nil {
			return nil, err
		}

		if specialist == nil {
			continue
		}

		response := domain.BookingResponse{
			ID:         booking.ID,
			UserID:     booking.UserID,
			TimeSlotID: booking.TimeSlotID,
			CreatedAt:  booking.CreatedAt,
			Specialist: specialist.Name,
			SlotTime:   timeSlot.Time,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
