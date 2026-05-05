package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/pkg/validator"
)

type CreateAppointmentInput struct {
	ServiceID string    `json:"service_id"`
	StartTime time.Time `json:"start_time"`
	Notes     string    `json:"notes"`
}

type AvailableSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type appointmentService struct {
	appointments repository.AppointmentRepository
	services     repository.ServiceRepository
	users        repository.UserRepository
	location     *time.Location
}

func NewAppointmentService(appointments repository.AppointmentRepository, services repository.ServiceRepository, users repository.UserRepository, location *time.Location) AppointmentService {
	if location == nil {
		location = time.UTC
	}
	return &appointmentService{
		appointments: appointments,
		services:     services,
		users:        users,
		location:     location,
	}
}

func (s *appointmentService) Create(ctx context.Context, clientID string, input CreateAppointmentInput) (*domain.Appointment, error) {
	var validationErrs validator.ValidationErrors
	if err := validator.ValidateRequired(input.ServiceID); err != nil {
		validationErrs.Add("service_id", err.Error())
	}
	if input.StartTime.IsZero() {
		validationErrs.Add("start_time", "is required")
	}
	if validationErrs.HasErrors() {
		return nil, validationErrs
	}

	service, err := s.services.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}
	if !service.IsActive || service.DeletedAt != nil {
		return nil, domain.ErrServiceInactive
	}
	if !input.StartTime.After(time.Now().UTC()) {
		return nil, domain.ErrPastTime
	}
	endTime := input.StartTime.Add(time.Duration(service.DurationMinutes) * time.Minute)

	appointment := &domain.Appointment{
		ClientID:  clientID,
		ServiceID: input.ServiceID,
		StartTime: input.StartTime.UTC(),
		EndTime:   endTime.UTC(),
		Status:    domain.AppointmentPending,
		Notes:     input.Notes,
	}

	if err := s.appointments.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		overlap, err := s.appointments.CheckOverlap(ctx, tx, input.ServiceID, appointment.StartTime, appointment.EndTime)
		if err != nil {
			return err
		}
		if overlap {
			return domain.ErrTimeSlotTaken
		}
		return s.appointments.Create(ctx, tx, appointment)
	}); err != nil {
		return nil, err
	}

	return s.appointments.GetByID(ctx, appointment.ID)
}

func (s *appointmentService) List(ctx context.Context, filter repository.AppointmentFilter) ([]domain.Appointment, repository.Pagination, error) {
	return s.appointments.List(ctx, filter)
}

func (s *appointmentService) ListMine(ctx context.Context, actorID string, role domain.Role, page, limit int) ([]domain.Appointment, repository.Pagination, error) {
	if role != domain.RoleClient && role != domain.RoleProvider {
		return nil, repository.Pagination{}, domain.ErrForbidden
	}
	return s.appointments.ListByActor(ctx, actorID, role, page, limit)
}

func (s *appointmentService) GetByID(ctx context.Context, actorID string, role domain.Role, id string) (*domain.Appointment, error) {
	appointment, err := s.appointments.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == domain.RoleAdmin || appointment.ClientID == actorID || appointment.ProviderID == actorID {
		return appointment, nil
	}
	return nil, domain.ErrForbidden
}

func (s *appointmentService) Confirm(ctx context.Context, actorID string, role domain.Role, id string) (*domain.Appointment, error) {
	appointment, err := s.GetByID(ctx, actorID, role, id)
	if err != nil {
		return nil, err
	}
	if role != domain.RoleAdmin && appointment.ProviderID != actorID {
		return nil, domain.ErrForbidden
	}
	if appointment.Status != domain.AppointmentPending {
		return nil, validator.ValidationErrors{{Field: "status", Error: "only pending appointments can be confirmed"}}
	}
	if err := s.appointments.UpdateStatus(ctx, id, domain.AppointmentConfirmed, "", time.Now().UTC()); err != nil {
		return nil, err
	}
	return s.appointments.GetByID(ctx, id)
}

func (s *appointmentService) Cancel(ctx context.Context, actorID string, role domain.Role, id string, reason string) (*domain.Appointment, error) {
	appointment, err := s.GetByID(ctx, actorID, role, id)
	if err != nil {
		return nil, err
	}
	if appointment.Status == domain.AppointmentCancelled || appointment.Status == domain.AppointmentCompleted {
		return nil, validator.ValidationErrors{{Field: "status", Error: "appointment cannot be cancelled in its current state"}}
	}
	if err := s.appointments.UpdateStatus(ctx, id, domain.AppointmentCancelled, reason, time.Now().UTC()); err != nil {
		return nil, err
	}
	return s.appointments.GetByID(ctx, id)
}

func (s *appointmentService) Complete(ctx context.Context, actorID string, role domain.Role, id string) (*domain.Appointment, error) {
	appointment, err := s.GetByID(ctx, actorID, role, id)
	if err != nil {
		return nil, err
	}
	if role != domain.RoleAdmin && appointment.ProviderID != actorID {
		return nil, domain.ErrForbidden
	}
	if appointment.Status != domain.AppointmentConfirmed {
		return nil, validator.ValidationErrors{{Field: "status", Error: "only confirmed appointments can be completed"}}
	}
	if err := s.appointments.UpdateStatus(ctx, id, domain.AppointmentCompleted, "", time.Now().UTC()); err != nil {
		return nil, err
	}
	return s.appointments.GetByID(ctx, id)
}

func (s *appointmentService) AvailableSlots(ctx context.Context, serviceID string, date time.Time) ([]AvailableSlot, error) {
	service, err := s.services.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	localDate := date.In(s.location)
	dayStart := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 9, 0, 0, 0, s.location)
	dayEnd := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 17, 0, 0, 0, s.location)
	fromDate := dayStart.UTC()
	toDate := dayEnd.UTC()
	appointments, _, err := s.appointments.List(ctx, repository.AppointmentFilter{
		Page:     1,
		Limit:    500,
		FromDate: &fromDate,
		ToDate:   &toDate,
	})
	if err != nil {
		return nil, err
	}

	var slots []AvailableSlot
	step := time.Duration(service.DurationMinutes) * time.Minute
	for current := dayStart; current.Add(step).Before(dayEnd) || current.Add(step).Equal(dayEnd); current = current.Add(step) {
		next := current.Add(step)
		available := true
		for _, appointment := range appointments {
			if appointment.ServiceID != serviceID {
				continue
			}
			if appointment.Status != domain.AppointmentPending && appointment.Status != domain.AppointmentConfirmed {
				continue
			}
			appointmentStart := appointment.StartTime.In(s.location)
			appointmentEnd := appointment.EndTime.In(s.location)
			if current.Before(appointmentEnd) && next.After(appointmentStart) {
				available = false
				break
			}
		}
		if available {
			slots = append(slots, AvailableSlot{
				StartTime: current.Format("15:04"),
				EndTime:   next.Format("15:04"),
			})
		}
	}
	return slots, nil
}
