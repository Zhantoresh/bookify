package usecase

import (
	"errors"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
	"github.com/bookify/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo *repository.UserRepository
}

func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Register(email, password string, role domain.Role) (*domain.User, error) {
	hashedPassword, err := u.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	if err := u.repo.CreateUser(user); err != nil {
		return nil, err
	}

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
		return "", err
	}

	if !u.ComparePassword(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return auth.GenerateToken(user.ID, string(user.Role))
}
