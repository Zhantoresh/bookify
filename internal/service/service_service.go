package service

import (
	"context"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/pkg/validator"
)

type CreateServiceInput struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
}

type UpdateServiceInput struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	IsActive        bool    `json:"is_active"`
}

type PatchServiceInput struct {
	Name            *string  `json:"name,omitempty"`
	Description     *string  `json:"description,omitempty"`
	Price           *float64 `json:"price,omitempty"`
	DurationMinutes *int     `json:"duration_minutes,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

type serviceService struct {
	services repository.ServiceRepository
	users    repository.UserRepository
}

func NewServiceService(services repository.ServiceRepository, users repository.UserRepository) ServiceService {
	return &serviceService{services: services, users: users}
}

func (s *serviceService) Create(ctx context.Context, providerID string, actorRole domain.Role, input CreateServiceInput) (*domain.Service, error) {
	if actorRole != domain.RoleProvider && actorRole != domain.RoleAdmin {
		return nil, domain.ErrForbidden
	}
	service := &domain.Service{
		ProviderID:      providerID,
		Name:            input.Name,
		Description:     input.Description,
		Price:           input.Price,
		DurationMinutes: input.DurationMinutes,
		IsActive:        true,
	}
	if err := validateServicePayload(service.Name, service.Price, service.DurationMinutes); err != nil {
		return nil, err
	}
	if err := s.services.Create(ctx, service); err != nil {
		return nil, err
	}
	return s.services.GetByID(ctx, service.ID)
}

func (s *serviceService) List(ctx context.Context, filter repository.ServiceFilter) ([]domain.Service, repository.Pagination, error) {
	filter.OnlyActive = true
	return s.services.List(ctx, filter)
}

func (s *serviceService) ListMine(ctx context.Context, providerID string, page, limit int) ([]domain.Service, repository.Pagination, error) {
	return s.services.List(ctx, repository.ServiceFilter{
		Page:       page,
		Limit:      limit,
		ProviderID: providerID,
	})
}

func (s *serviceService) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	return s.services.GetByID(ctx, id)
}

func (s *serviceService) Update(ctx context.Context, actorID string, actorRole domain.Role, id string, input UpdateServiceInput) (*domain.Service, error) {
	service, err := s.services.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureServiceAccess(actorID, actorRole, service.ProviderID); err != nil {
		return nil, err
	}
	service.Name = input.Name
	service.Description = input.Description
	service.Price = input.Price
	service.DurationMinutes = input.DurationMinutes
	service.IsActive = input.IsActive
	if !service.IsActive && service.DeletedAt == nil {
		now := time.Now().UTC()
		service.DeletedAt = &now
	}
	if err := validateServicePayload(service.Name, service.Price, service.DurationMinutes); err != nil {
		return nil, err
	}
	if err := s.services.Update(ctx, service); err != nil {
		return nil, err
	}
	return s.services.GetByID(ctx, id)
}

func (s *serviceService) Patch(ctx context.Context, actorID string, actorRole domain.Role, id string, input PatchServiceInput) (*domain.Service, error) {
	service, err := s.services.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := ensureServiceAccess(actorID, actorRole, service.ProviderID); err != nil {
		return nil, err
	}
	if input.Name != nil {
		service.Name = *input.Name
	}
	if input.Description != nil {
		service.Description = *input.Description
	}
	if input.Price != nil {
		service.Price = *input.Price
	}
	if input.DurationMinutes != nil {
		service.DurationMinutes = *input.DurationMinutes
	}
	if input.IsActive != nil {
		service.IsActive = *input.IsActive
		if !*input.IsActive && service.DeletedAt == nil {
			now := time.Now().UTC()
			service.DeletedAt = &now
		}
	}
	if err := validateServicePayload(service.Name, service.Price, service.DurationMinutes); err != nil {
		return nil, err
	}
	if err := s.services.Update(ctx, service); err != nil {
		return nil, err
	}
	return s.services.GetByID(ctx, id)
}

func (s *serviceService) Delete(ctx context.Context, actorID string, actorRole domain.Role, id string) error {
	service, err := s.services.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := ensureServiceAccess(actorID, actorRole, service.ProviderID); err != nil {
		return err
	}
	hasFuture, err := s.services.HasFutureAppointments(ctx, id, time.Now().UTC())
	if err != nil {
		return err
	}
	if hasFuture {
		return domain.ErrHasFutureBookings
	}
	service.IsActive = false
	now := time.Now().UTC()
	service.DeletedAt = &now
	return s.services.Update(ctx, service)
}

func ensureServiceAccess(actorID string, actorRole domain.Role, providerID string) error {
	if actorRole == domain.RoleAdmin {
		return nil
	}
	if actorRole == domain.RoleProvider && actorID == providerID {
		return nil
	}
	return domain.ErrForbidden
}

func validateServicePayload(name string, price float64, durationMinutes int) error {
	var validationErrs validator.ValidationErrors
	if name == "" {
		validationErrs.Add("name", "is required")
	}
	if price < 0 {
		validationErrs.Add("price", "must be greater than or equal to 0")
	}
	if durationMinutes <= 0 {
		validationErrs.Add("duration_minutes", "must be greater than 0")
	}
	if validationErrs.HasErrors() {
		return validationErrs
	}
	return nil
}
