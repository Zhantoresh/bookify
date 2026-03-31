package service

import (
	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/notification"
	"github.com/bookify/internal/repository"
)

type BookingService struct {
	bookingRepo    *repository.BookingRepository
	timeSlotRepo   *repository.TimeSlotRepository
	specialistRepo *repository.SpecialistRepository
	userRepo       *repository.UserRepository
	notifier       notification.Notifier
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	timeSlotRepo *repository.TimeSlotRepository,
	specialistRepo *repository.SpecialistRepository,
	userRepo *repository.UserRepository,
	notifier notification.Notifier,
) *BookingService {
	if notifier == nil {
		notifier = notification.NewNoopNotifier()
	}

	return &BookingService{
		bookingRepo:    bookingRepo,
		timeSlotRepo:   timeSlotRepo,
		specialistRepo: specialistRepo,
		userRepo:       userRepo,
		notifier:       notifier,
	}
}

func (s *BookingService) CreateBooking(userID, timeSlotID int) (*domain.Booking, error) {
	// Check if time slot exists
	timeSlot, err := s.timeSlotRepo.GetByID(timeSlotID)
	if err != nil {
		return nil, err
	}

	if timeSlot == nil {
		return nil, ErrTimeSlotNotFound
	}

	// Check if slot is already booked
	if timeSlot.IsBooked {
		return nil, ErrBookingAlreadyBooked
	}

	// Create booking
	booking, err := s.bookingRepo.Create(userID, timeSlotID)
	if err != nil {
		return nil, err
	}

	// Mark slot as booked
	err = s.timeSlotRepo.MarkAsBooked(timeSlotID)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) CreateBookingWithDetails(userID, timeSlotID int) (*domain.CreateBookingResponse, error) {
	// Create booking
	booking, err := s.CreateBooking(userID, timeSlotID)
	if err != nil {
		return nil, err
	}

	// Get time slot info
	timeSlot, err := s.timeSlotRepo.GetByID(booking.TimeSlotID)
	if err != nil {
		return nil, err
	}

	if timeSlot == nil {
		return nil, ErrTimeSlotNotFound
	}

	// Get specialist info
	specialist, err := s.specialistRepo.GetByID(timeSlot.UserID)
	if err != nil {
		return nil, err
	}

	if specialist == nil {
		return nil, ErrSpecialistNotFound
	}

	if s.userRepo != nil {
		if user, userErr := s.userRepo.GetByID(userID); userErr == nil && user != nil {
			s.notifier.Notify(notification.BuildBookingCreatedMessage(user.Email, user.Name, specialist.Name, timeSlot.Time))
		}
	}

	return &domain.CreateBookingResponse{
		ID:         booking.ID,
		Specialist: specialist.Name,
		Time:       timeSlot.Time,
		Status:     booking.Status,
	}, nil
}

func (s *BookingService) CancelBooking(userID, bookingID int) error {
	// Get booking
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking == nil {
		return ErrBookingNotFound
	}

	// Check ownership
	if booking.UserID != userID {
		return ErrForbidden
	}

	// Check if already cancelled
	if booking.Status == "CANCELLED" {
		return ErrBookingAlreadyCancelled
	}

	// Cancel booking
	err = s.bookingRepo.CancelByID(bookingID)
	if err != nil {
		return err
	}

	// Free up the slot
	err = s.timeSlotRepo.MarkAsUnbooked(booking.TimeSlotID)
	if err != nil {
		return err
	}

	return nil
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

		specialist, err := s.specialistRepo.GetByID(timeSlot.UserID)
		if err != nil {
			return nil, err
		}

		if specialist == nil {
			continue
		}

		response := domain.BookingResponse{
			ID:         booking.ID,
			Specialist: specialist.Name,
			Time:       timeSlot.Time,
			Status:     booking.Status,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
