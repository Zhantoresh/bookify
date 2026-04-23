package service

import (
	"context"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/pkg/validator"
)

type adminService struct {
	users        repository.UserRepository
	services     repository.ServiceRepository
	appointments repository.AppointmentRepository
}

func NewAdminService(users repository.UserRepository, services repository.ServiceRepository, appointments repository.AppointmentRepository) AdminService {
	return &adminService{
		users:        users,
		services:     services,
		appointments: appointments,
	}
}

func (s *adminService) Dashboard(ctx context.Context) (*Dashboard, error) {
	userCounts, err := s.users.CountByRole(ctx)
	if err != nil {
		return nil, err
	}
	totalServices, err := s.services.Count(ctx, false)
	if err != nil {
		return nil, err
	}
	activeServices, err := s.services.Count(ctx, true)
	if err != nil {
		return nil, err
	}
	appointmentCounts, err := s.appointments.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}
	recentUsers, err := s.users.List(ctx, repository.UserFilter{})
	if err != nil {
		return nil, err
	}
	if len(recentUsers) > 5 {
		recentUsers = recentUsers[:5]
	}
	recentServices, _, err := s.services.List(ctx, repository.ServiceFilter{Page: 1, Limit: 5})
	if err != nil {
		return nil, err
	}
	recentAppointments, _, err := s.appointments.List(ctx, repository.AppointmentFilter{Page: 1, Limit: 5})
	if err != nil {
		return nil, err
	}

	totalUsers := userCounts[domain.RoleAdmin] + userCounts[domain.RoleClient] + userCounts[domain.RoleProvider]
	totalAppointments := appointmentCounts[domain.AppointmentPending] + appointmentCounts[domain.AppointmentConfirmed] + appointmentCounts[domain.AppointmentCancelled] + appointmentCounts[domain.AppointmentCompleted]

	return &Dashboard{
		Summary: DashboardSummary{
			TotalUsers:            totalUsers,
			TotalClients:          userCounts[domain.RoleClient],
			TotalProviders:        userCounts[domain.RoleProvider],
			TotalAdmins:           userCounts[domain.RoleAdmin],
			TotalServices:         totalServices,
			ActiveServices:        activeServices,
			TotalAppointments:     totalAppointments,
			PendingAppointments:   appointmentCounts[domain.AppointmentPending],
			ConfirmedAppointments: appointmentCounts[domain.AppointmentConfirmed],
			CancelledAppointments: appointmentCounts[domain.AppointmentCancelled],
			CompletedAppointments: appointmentCounts[domain.AppointmentCompleted],
		},
		RecentUsers:        recentUsers,
		RecentServices:     recentServices,
		RecentAppointments: recentAppointments,
	}, nil
}

func (s *adminService) ListUsers(ctx context.Context, filter repository.UserFilter) ([]domain.User, error) {
	return s.users.List(ctx, filter)
}

func (s *adminService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *adminService) UpdateUserRole(ctx context.Context, actorID, targetUserID string, role domain.Role) (*domain.User, error) {
	if actorID == targetUserID {
		return nil, validator.ValidationErrors{{Field: "user_id", Error: "admin cannot change own role"}}
	}
	if role != domain.RoleAdmin && role != domain.RoleClient && role != domain.RoleProvider {
		return nil, validator.ValidationErrors{{Field: "role", Error: "must be one of: admin, client, provider"}}
	}
	target, err := s.users.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	if err := s.users.UpdateRole(ctx, target.ID, role); err != nil {
		return nil, err
	}
	return s.users.GetByID(ctx, target.ID)
}

func (s *adminService) DeleteUser(ctx context.Context, actorID, targetUserID string) error {
	if actorID == targetUserID {
		return validator.ValidationErrors{{Field: "user_id", Error: "admin cannot delete own account"}}
	}
	return s.users.Delete(ctx, targetUserID)
}
