package service

import (
	"errors"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type TimeSlotService struct {
	timeSlotRepo *repository.TimeSlotRepository
}

func NewTimeSlotService(timeSlotRepo *repository.TimeSlotRepository) *TimeSlotService {
	return &TimeSlotService{
		timeSlotRepo: timeSlotRepo,
	}
}

// CreateTimeSlot creates a new time slot for the authenticated specialist
func (s *TimeSlotService) CreateTimeSlot(userID int, slotTime time.Time) (*domain.TimeSlot, error) {
	return s.timeSlotRepo.CreateTimeSlot(userID, slotTime)
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
		return errors.New("time slot not found")
	}

	if slot.UserID != userID {
		return errors.New("forbidden")
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
		return errors.New("time slot not found")
	}

	if slot.UserID != userID {
		return errors.New("forbidden")
	}

	return s.timeSlotRepo.DeleteTimeSlot(slotID)
}
