package service

import (
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/notification"
	"github.com/bookify/internal/repository"
)

type TimeSlotService struct {
	timeSlotRepo *repository.TimeSlotRepository
	userRepo     *repository.UserRepository
	notifier     notification.Notifier
}

func NewTimeSlotService(
	timeSlotRepo *repository.TimeSlotRepository,
	userRepo *repository.UserRepository,
	notifier notification.Notifier,
) *TimeSlotService {
	if notifier == nil {
		notifier = notification.NewNoopNotifier()
	}

	return &TimeSlotService{
		timeSlotRepo: timeSlotRepo,
		userRepo:     userRepo,
		notifier:     notifier,
	}
}

// CreateTimeSlot creates a new time slot for the authenticated specialist
func (s *TimeSlotService) CreateTimeSlot(userID int, slotTime time.Time) (*domain.TimeSlot, error) {
	slot, err := s.timeSlotRepo.CreateTimeSlot(userID, slotTime)
	if err != nil {
		return nil, err
	}

	if s.userRepo != nil {
		if user, userErr := s.userRepo.GetByID(userID); userErr == nil && user != nil {
			s.notifier.Notify(notification.BuildTimeSlotCreatedMessage(user.Email, user.Name, slot.Time))
		}
	}

	return slot, nil
}

// GetMyTimeSlots returns all time slots for the authenticated specialist
func (s *TimeSlotService) GetMyTimeSlots(userID int) ([]domain.TimeSlot, error) {
	return s.timeSlotRepo.GetByUserID(userID)
}

// UpdateTimeSlot updates a time slot, ensuring the user owns it
func (s *TimeSlotService) UpdateTimeSlot(slotID, userID int, slotTime time.Time) error {
	slot, err := s.timeSlotRepo.GetByID(slotID)
	if err != nil {
		return err
	}

	if slot == nil {
		return ErrTimeSlotNotFound
	}

	if slot.UserID != userID {
		return ErrForbidden
	}

	return s.timeSlotRepo.UpdateTimeSlot(slotID, slotTime)
}

// DeleteTimeSlot deletes a time slot, ensuring the user owns it
func (s *TimeSlotService) DeleteTimeSlot(slotID, userID int) error {
	slot, err := s.timeSlotRepo.GetByID(slotID)
	if err != nil {
		return err
	}

	if slot == nil {
		return ErrTimeSlotNotFound
	}

	if slot.UserID != userID {
		return ErrForbidden
	}

	return s.timeSlotRepo.DeleteTimeSlot(slotID)
}
