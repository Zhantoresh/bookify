package service

import (
	"errors"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type SpecialistService struct {
	specialistRepo *repository.SpecialistRepository
	timeSlotRepo   *repository.TimeSlotRepository
}

func NewSpecialistService(
	specialistRepo *repository.SpecialistRepository,
	timeSlotRepo *repository.TimeSlotRepository,
) *SpecialistService {
	return &SpecialistService{
		specialistRepo: specialistRepo,
		timeSlotRepo:   timeSlotRepo,
	}
}

func (s *SpecialistService) GetAllSpecialists() ([]domain.Specialist, error) {
	return s.specialistRepo.GetAll()
}

type SpecialistWithSlots struct {
	Specialist domain.Specialist `json:"specialist"`
	TimeSlots  []domain.TimeSlot `json:"time_slots"`
}

func (s *SpecialistService) GetSpecialistWithSlots(id int) (*SpecialistWithSlots, error) {
	specialist, err := s.specialistRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if specialist == nil {
		return nil, errors.New("specialist not found")
	}

	// Get time slots for this specialist (user)
	slots, err := s.timeSlotRepo.GetByUserID(id)
	if err != nil {
		return nil, err
	}

	return &SpecialistWithSlots{
		Specialist: *specialist,
		TimeSlots:  slots,
	}, nil
}
