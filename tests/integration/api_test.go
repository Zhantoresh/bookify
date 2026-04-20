package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	appservice "github.com/bookify/internal/service"
	authsvc "github.com/bookify/internal/service/auth"
	httptransport "github.com/bookify/internal/transport/http"
	"github.com/bookify/internal/worker"
	"github.com/bookify/pkg/logger"
)

type inMemoryUserRepo struct {
	byID    map[string]*domain.User
	byEmail map[string]*domain.User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		byID:    map[string]*domain.User{},
		byEmail: map[string]*domain.User{},
	}
}

func (r *inMemoryUserRepo) Create(ctx context.Context, user *domain.User) error {
	user.ID = "user-" + user.Email
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *inMemoryUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepo) List(ctx context.Context, filter repository.UserFilter) ([]domain.User, error) {
	return nil, nil
}

type noopServiceRepo struct{}

func (n *noopServiceRepo) Create(ctx context.Context, service *domain.Service) error { return nil }
func (n *noopServiceRepo) List(ctx context.Context, filter repository.ServiceFilter) ([]domain.Service, repository.Pagination, error) {
	return nil, repository.Pagination{}, nil
}
func (n *noopServiceRepo) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	return nil, domain.ErrNotFound
}
func (n *noopServiceRepo) Update(ctx context.Context, service *domain.Service) error { return nil }
func (n *noopServiceRepo) HasFutureAppointments(ctx context.Context, serviceID string, now time.Time) (bool, error) {
	return false, nil
}

type noopAppointmentRepo struct{}

func (n *noopAppointmentRepo) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	return nil
}
func (n *noopAppointmentRepo) CheckOverlap(ctx context.Context, tx *sql.Tx, serviceID string, start, end time.Time) (bool, error) {
	return false, nil
}
func (n *noopAppointmentRepo) Create(ctx context.Context, tx *sql.Tx, appointment *domain.Appointment) error {
	return nil
}
func (n *noopAppointmentRepo) List(ctx context.Context, filter repository.AppointmentFilter) ([]domain.Appointment, repository.Pagination, error) {
	return nil, repository.Pagination{}, nil
}
func (n *noopAppointmentRepo) ListByActor(ctx context.Context, actorID string, role domain.Role, page, limit int) ([]domain.Appointment, repository.Pagination, error) {
	return nil, repository.Pagination{}, nil
}
func (n *noopAppointmentRepo) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	return nil, domain.ErrNotFound
}
func (n *noopAppointmentRepo) UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus, reason string, changedAt time.Time) error {
	return nil
}
func (n *noopAppointmentRepo) GetAppointmentsByDateRange(ctx context.Context, start, end time.Time) ([]domain.Appointment, error) {
	return nil, nil
}

func setupTestRouter() http.Handler {
	userRepo := newInMemoryUserRepo()
	jwt := authsvc.NewJWTService("test-secret", 24*time.Hour)
	authService := appservice.NewAuthService(userRepo, jwt)
	userService := appservice.NewUserService(userRepo)
	serviceService := appservice.NewServiceService(&noopServiceRepo{}, userRepo)
	appointmentService := appservice.NewAppointmentService(&noopAppointmentRepo{}, &noopServiceRepo{}, userRepo)
	log := logger.New("error")
	pool := worker.NewWorkerPool(1, 5, log)
	return httptransport.NewServer(authService, userService, serviceService, appointmentService, jwt, pool, log)
}

func TestRegisterUser(t *testing.T) {
	router := setupTestRouter()
	body := map[string]string{
		"email":     "test@example.com",
		"password":  "TestPass123",
		"full_name": "Test User",
		"role":      "client",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestLoginUser(t *testing.T) {
	router := setupTestRouter()

	registerBody := map[string]string{
		"email":     "login@example.com",
		"password":  "LoginPass123",
		"full_name": "Login User",
		"role":      "client",
	}
	payload, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	loginBody := map[string]string{
		"email":    "login@example.com",
		"password": "LoginPass123",
	}
	payload, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var response map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	if response["token"] == "" {
		t.Fatal("expected token in response")
	}
}

func TestGetServicesMyUnauthorized(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/my", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestValidateToken(t *testing.T) {
	router := setupTestRouter()

	registerBody := map[string]string{
		"email":     "validate@example.com",
		"password":  "Validate123",
		"full_name": "Validate User",
		"role":      "client",
	}
	payload, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	loginBody := map[string]string{
		"email":    "validate@example.com",
		"password": "Validate123",
	}
	payload, _ = json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var response map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	token, _ := response["token"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRegisterValidationErrorFormat(t *testing.T) {
	router := setupTestRouter()
	body := map[string]string{
		"email":    "bad",
		"password": "123",
		"role":     "bad",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	var response map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	if response["error"] == nil {
		t.Fatal("expected error envelope")
	}
}
