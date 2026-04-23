package unit

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	appservice "github.com/bookify/internal/service"
	"github.com/bookify/pkg/validator"
)

type fakeAppointmentRepo struct {
	overlap     bool
	created     *domain.Appointment
	stored      *domain.Appointment
	updateError error
}

func (f *fakeAppointmentRepo) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	return fn(ctx, nil)
}

func (f *fakeAppointmentRepo) CheckOverlap(ctx context.Context, tx *sql.Tx, serviceID string, start, end time.Time) (bool, error) {
	return f.overlap, nil
}

func (f *fakeAppointmentRepo) Create(ctx context.Context, tx *sql.Tx, appointment *domain.Appointment) error {
	appointment.ID = "appointment-1"
	appointment.CreatedAt = time.Now().UTC()
	f.created = appointment
	f.stored = appointment
	return nil
}

func (f *fakeAppointmentRepo) List(ctx context.Context, filter repository.AppointmentFilter) ([]domain.Appointment, repository.Pagination, error) {
	return nil, repository.Pagination{}, nil
}

func (f *fakeAppointmentRepo) ListByActor(ctx context.Context, actorID string, role domain.Role, page, limit int) ([]domain.Appointment, repository.Pagination, error) {
	return nil, repository.Pagination{}, nil
}

func (f *fakeAppointmentRepo) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	if f.stored == nil {
		return nil, domain.ErrNotFound
	}
	return f.stored, nil
}

func (f *fakeAppointmentRepo) UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus, reason string, changedAt time.Time) error {
	return f.updateError
}

func (f *fakeAppointmentRepo) GetAppointmentsByDateRange(ctx context.Context, start, end time.Time) ([]domain.Appointment, error) {
	return nil, nil
}

func (f *fakeAppointmentRepo) CountByStatus(ctx context.Context) (map[domain.AppointmentStatus]int, error) {
	return map[domain.AppointmentStatus]int{}, nil
}

type fakeServiceRepo struct {
	service *domain.Service
}

func (f *fakeServiceRepo) Create(ctx context.Context, service *domain.Service) error {
	return nil
}

func (f *fakeServiceRepo) List(ctx context.Context, filter repository.ServiceFilter) ([]domain.Service, repository.Pagination, error) {
	return nil, repository.Pagination{}, nil
}

func (f *fakeServiceRepo) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	if f.service == nil {
		return nil, domain.ErrNotFound
	}
	return f.service, nil
}

func (f *fakeServiceRepo) Update(ctx context.Context, service *domain.Service) error {
	return nil
}

func (f *fakeServiceRepo) HasFutureAppointments(ctx context.Context, serviceID string, now time.Time) (bool, error) {
	return false, nil
}

func (f *fakeServiceRepo) Count(ctx context.Context, onlyActive bool) (int, error) {
	return 0, nil
}

type fakeUserRepo struct{}

func (f *fakeUserRepo) Create(ctx context.Context, user *domain.User) error { return nil }
func (f *fakeUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeUserRepo) List(ctx context.Context, filter repository.UserFilter) ([]domain.User, error) {
	return nil, nil
}
func (f *fakeUserRepo) CountByRole(ctx context.Context) (map[domain.Role]int, error) {
	return map[domain.Role]int{}, nil
}
func (f *fakeUserRepo) UpdateRole(ctx context.Context, id string, role domain.Role) error {
	return nil
}
func (f *fakeUserRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func TestCreateAppointmentSuccess(t *testing.T) {
	repo := &fakeAppointmentRepo{}
	serviceRepo := &fakeServiceRepo{
		service: &domain.Service{
			ID:              "service-1",
			DurationMinutes: 30,
			IsActive:        true,
		},
	}
	svc := appservice.NewAppointmentService(repo, serviceRepo, &fakeUserRepo{})

	result, err := svc.Create(context.Background(), "client-1", appservice.CreateAppointmentInput{
		ServiceID: "service-1",
		StartTime: time.Now().UTC().Add(24 * time.Hour),
		Notes:     "note",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.Status != domain.AppointmentPending {
		t.Fatalf("expected pending appointment, got %+v", result)
	}
}

func TestCreateAppointmentOverlap(t *testing.T) {
	repo := &fakeAppointmentRepo{overlap: true}
	serviceRepo := &fakeServiceRepo{
		service: &domain.Service{
			ID:              "service-1",
			DurationMinutes: 30,
			IsActive:        true,
		},
	}
	svc := appservice.NewAppointmentService(repo, serviceRepo, &fakeUserRepo{})

	result, err := svc.Create(context.Background(), "client-1", appservice.CreateAppointmentInput{
		ServiceID: "service-1",
		StartTime: time.Now().UTC().Add(24 * time.Hour),
	})
	if !errors.Is(err, domain.ErrTimeSlotTaken) {
		t.Fatalf("expected ErrTimeSlotTaken, got %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil appointment, got %+v", result)
	}
}

func TestCompleteAppointmentInvalidStatus(t *testing.T) {
	repo := &fakeAppointmentRepo{
		stored: &domain.Appointment{
			ID:         "appointment-1",
			ProviderID: "provider-1",
			Status:     domain.AppointmentPending,
		},
	}
	svc := appservice.NewAppointmentService(repo, &fakeServiceRepo{}, &fakeUserRepo{})

	result, err := svc.Complete(context.Background(), "provider-1", domain.RoleProvider, "appointment-1")
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		t.Fatalf("expected validation errors, got %v", err)
	}
}
