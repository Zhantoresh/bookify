package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	appservice "github.com/bookify/internal/service"
	authsvc "github.com/bookify/internal/service/auth"
	"github.com/bookify/pkg/validator"
)

func TestHashPassword(t *testing.T) {
	hash, err := authsvc.HashPassword("SecurePass123")
	if err != nil {
		t.Fatalf("expected hash generation to succeed: %v", err)
	}
	if !authsvc.CheckPasswordHash("SecurePass123", hash) {
		t.Fatal("expected password hash to validate")
	}
}

type authUserRepo struct {
	user *domain.User
}

func (a *authUserRepo) Create(ctx context.Context, user *domain.User) error {
	a.user = user
	user.ID = "user-1"
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	return nil
}

func (a *authUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return a.user, nil
}

func (a *authUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if a.user == nil || a.user.Email != email {
		return nil, domain.ErrNotFound
	}
	return a.user, nil
}

func (a *authUserRepo) List(ctx context.Context, filter repository.UserFilter) ([]domain.User, error) {
	return nil, nil
}

func TestRegisterValidation(t *testing.T) {
	svc := appservice.NewAuthService(&authUserRepo{}, authsvc.NewJWTService("secret", time.Hour))
	_, err := svc.Register(context.Background(), appservice.RegisterInput{
		Email:    "bad-email",
		Password: "123",
		FullName: "",
		Role:     "wrong",
	})
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		t.Fatalf("expected validation errors, got %v", err)
	}
}
