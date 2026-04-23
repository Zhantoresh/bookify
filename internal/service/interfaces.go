package service

import (
	"context"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
	ValidateToken(token string) (*Claims, error)
}

type UserService interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, filter repository.UserFilter) ([]domain.User, error)
}

type AdminService interface {
	Dashboard(ctx context.Context) (*Dashboard, error)
	ListUsers(ctx context.Context, filter repository.UserFilter) ([]domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateUserRole(ctx context.Context, actorID, targetUserID string, role domain.Role) (*domain.User, error)
	DeleteUser(ctx context.Context, actorID, targetUserID string) error
}

type ServiceService interface {
	Create(ctx context.Context, providerID string, actorRole domain.Role, input CreateServiceInput) (*domain.Service, error)
	List(ctx context.Context, filter repository.ServiceFilter) ([]domain.Service, repository.Pagination, error)
	ListMine(ctx context.Context, providerID string, page, limit int) ([]domain.Service, repository.Pagination, error)
	GetByID(ctx context.Context, id string) (*domain.Service, error)
	Update(ctx context.Context, actorID string, actorRole domain.Role, id string, input UpdateServiceInput) (*domain.Service, error)
	Patch(ctx context.Context, actorID string, actorRole domain.Role, id string, input PatchServiceInput) (*domain.Service, error)
	Delete(ctx context.Context, actorID string, actorRole domain.Role, id string) error
}

type AppointmentService interface {
	Create(ctx context.Context, clientID string, input CreateAppointmentInput) (*domain.Appointment, error)
	List(ctx context.Context, filter repository.AppointmentFilter) ([]domain.Appointment, repository.Pagination, error)
	ListMine(ctx context.Context, actorID string, role domain.Role, page, limit int) ([]domain.Appointment, repository.Pagination, error)
	GetByID(ctx context.Context, actorID string, role domain.Role, id string) (*domain.Appointment, error)
	Confirm(ctx context.Context, actorID string, role domain.Role, id string) (*domain.Appointment, error)
	Cancel(ctx context.Context, actorID string, role domain.Role, id string, reason string) (*domain.Appointment, error)
	Complete(ctx context.Context, actorID string, role domain.Role, id string) (*domain.Appointment, error)
	AvailableSlots(ctx context.Context, serviceID string, date time.Time) ([]AvailableSlot, error)
}

type Dashboard struct {
	Summary           DashboardSummary       `json:"summary"`
	RecentUsers       []domain.User          `json:"recent_users"`
	RecentServices    []domain.Service       `json:"recent_services"`
	RecentAppointments []domain.Appointment  `json:"recent_appointments"`
}

type DashboardSummary struct {
	TotalUsers               int `json:"total_users"`
	TotalClients             int `json:"total_clients"`
	TotalProviders           int `json:"total_providers"`
	TotalAdmins              int `json:"total_admins"`
	TotalServices            int `json:"total_services"`
	ActiveServices           int `json:"active_services"`
	TotalAppointments        int `json:"total_appointments"`
	PendingAppointments      int `json:"pending_appointments"`
	ConfirmedAppointments    int `json:"confirmed_appointments"`
	CancelledAppointments    int `json:"cancelled_appointments"`
	CompletedAppointments    int `json:"completed_appointments"`
}
