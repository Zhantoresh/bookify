package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/bookify/internal/domain"
)

type UserFilter struct {
	Role domain.Role
}

type ServiceFilter struct {
	Page       int
	Limit      int
	ProviderID string
	MinPrice   *float64
	MaxPrice   *float64
	Search     string
	OnlyActive bool
}

type AppointmentFilter struct {
	Page       int
	Limit      int
	Status     string
	FromDate   *time.Time
	ToDate     *time.Time
	ProviderID string
	ClientID   string
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, filter UserFilter) ([]domain.User, error)
	CountByRole(ctx context.Context) (map[domain.Role]int, error)
	UpdateRole(ctx context.Context, id string, role domain.Role) error
	Delete(ctx context.Context, id string) error
}

type ServiceRepository interface {
	Create(ctx context.Context, service *domain.Service) error
	List(ctx context.Context, filter ServiceFilter) ([]domain.Service, Pagination, error)
	GetByID(ctx context.Context, id string) (*domain.Service, error)
	Update(ctx context.Context, service *domain.Service) error
	HasFutureAppointments(ctx context.Context, serviceID string, now time.Time) (bool, error)
	Count(ctx context.Context, onlyActive bool) (int, error)
}

type AppointmentRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
	CheckOverlap(ctx context.Context, tx *sql.Tx, serviceID string, start, end time.Time) (bool, error)
	Create(ctx context.Context, tx *sql.Tx, appointment *domain.Appointment) error
	List(ctx context.Context, filter AppointmentFilter) ([]domain.Appointment, Pagination, error)
	ListByActor(ctx context.Context, actorID string, role domain.Role, page, limit int) ([]domain.Appointment, Pagination, error)
	GetByID(ctx context.Context, id string) (*domain.Appointment, error)
	UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus, reason string, changedAt time.Time) error
	GetAppointmentsByDateRange(ctx context.Context, start, end time.Time) ([]domain.Appointment, error)
	CountByStatus(ctx context.Context) (map[domain.AppointmentStatus]int, error)
}
