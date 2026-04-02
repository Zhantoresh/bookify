package service

import (
	"log/slog"

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
	logger         *slog.Logger
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	timeSlotRepo *repository.TimeSlotRepository,
	specialistRepo *repository.SpecialistRepository,
	userRepo *repository.UserRepository,
	notifier notification.Notifier,
	logger *slog.Logger,
) *BookingService {
	if notifier == nil {
		notifier = notification.NewNoopNotifier()
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &BookingService{
		bookingRepo:    bookingRepo,
		timeSlotRepo:   timeSlotRepo,
		specialistRepo: specialistRepo,
		userRepo:       userRepo,
		notifier:       notifier,
		logger:         logger,
	}
}

func (s *BookingService) CreateBooking(userID, timeSlotID int) (*domain.Booking, error) {
	timeSlot, err := s.timeSlotRepo.GetByID(timeSlotID)
	if err != nil {
		return nil, err
	}

	if timeSlot == nil {
		return nil, ErrTimeSlotNotFound
	}

	if timeSlot.IsBooked {
		return nil, ErrBookingAlreadyBooked
	}

	booking, err := s.bookingRepo.Create(userID, timeSlotID)
	if err != nil {
		return nil, err
	}

	err = s.timeSlotRepo.MarkAsBooked(timeSlotID)
	if err != nil {
		return nil, err
	}

	s.logger.Info("booking created",
		"user_id", userID,
		"booking_id", booking.ID,
		"time_slot_id", timeSlotID,
	)

	return booking, nil
}

func (s *BookingService) CreateBookingWithDetails(userID, timeSlotID int) (*domain.CreateBookingResponse, error) {
	booking, err := s.CreateBooking(userID, timeSlotID)
	if err != nil {
		return nil, err
	}

	timeSlot, err := s.timeSlotRepo.GetByID(booking.TimeSlotID)
	if err != nil {
		return nil, err
	}

	if timeSlot == nil {
		return nil, ErrTimeSlotNotFound
	}

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
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking == nil {
		return ErrBookingNotFound
	}

	if booking.UserID != userID {
		return ErrForbidden
	}

	if booking.Status == "CANCELLED" {
		return ErrBookingAlreadyCancelled
	}

	err = s.bookingRepo.CancelByID(bookingID)
	if err != nil {
		return err
	}

	err = s.timeSlotRepo.MarkAsUnbooked(booking.TimeSlotID)
	if err != nil {
		return err
	}

	s.logger.Info("booking cancelled",
		"user_id", userID,
		"booking_id", bookingID,
		"time_slot_id", booking.TimeSlotID,
	)

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