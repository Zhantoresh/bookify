package service

import (
	"context"
	"errors"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	authsvc "github.com/bookify/internal/service/auth"
	"github.com/bookify/pkg/validator"
)

type Claims = authsvc.Claims

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
}

type authService struct {
	users repository.UserRepository
	jwt   *authsvc.JWTService
}

func NewAuthService(users repository.UserRepository, jwt *authsvc.JWTService) AuthService {
	return &authService{users: users, jwt: jwt}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	var validationErrs validator.ValidationErrors
	if err := validator.ValidateEmail(input.Email); err != nil {
		validationErrs.Add("email", err.Error())
	}
	if err := validator.ValidatePassword(input.Password); err != nil {
		validationErrs.Add("password", err.Error())
	}
	if err := validator.ValidateRequired(input.FullName); err != nil {
		validationErrs.Add("full_name", err.Error())
	}
	if err := validator.ValidatePhone(input.Phone); err != nil {
		validationErrs.Add("phone", err.Error())
	}

	role := domain.Role(input.Role)
	if role != domain.RoleClient && role != domain.RoleProvider {
		validationErrs.Add("role", "must be either client or provider")
	}
	if validationErrs.HasErrors() {
		return nil, validationErrs
	}

	hash, err := authsvc.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: hash,
		FullName:     input.FullName,
		Role:         role,
		Phone:        input.Phone,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	var validationErrs validator.ValidationErrors
	if err := validator.ValidateEmail(email); err != nil {
		validationErrs.Add("email", err.Error())
	}
	if err := validator.ValidateRequired(password); err != nil {
		validationErrs.Add("password", err.Error())
	}
	if validationErrs.HasErrors() {
		return "", nil, validationErrs
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", nil, domain.ErrInvalidCredentials
		}
		return "", nil, err
	}
	if !authsvc.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, domain.ErrInvalidCredentials
	}
	token, err := s.jwt.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *authService) ValidateToken(token string) (*Claims, error) {
	return s.jwt.ValidateToken(token)
}
