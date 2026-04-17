package usecase

import (
	"errors"

	"log/slog"
	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo   *repository.UserRepository
	logger *slog.Logger 
}

func NewUserUsecase(repo *repository.UserRepository, logger *slog.Logger) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *UserUsecase) Register(email, password, name string, role domain.Role) (*domain.User, error) {
	hashedPassword, err := u.HashPassword(password)
	if err != nil {
		u.logger.Error("failed to hash password during registration", "email", email, "error", err)
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := u.repo.CreateUser(user); err != nil {
		u.logger.Error("failed to create user in database", "email", email, "error", err)
		return nil, err
	}

	u.logger.Info("user registered successfully", "email", email, "role", string(role))
	return user, nil
}

func (u *UserUsecase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserUsecase) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *UserUsecase) Login(email, password string) (string, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		u.logger.Warn("failed login attempt: user not found", "email", email)
		return "", err
	}

	if !u.ComparePassword(password, user.PasswordHash) {
		u.logger.Warn("failed login attempt: invalid password", "email", email)
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		u.logger.Error("failed to generate token during login", "user_id", user.ID, "error", err)
		return "", err
	}

	u.logger.Info("user logged in", "user_id", user.ID, "role", string(user.Role))

	return token, nil
}
