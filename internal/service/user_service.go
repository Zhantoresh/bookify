package service

import (
	"context"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type userService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context, filter repository.UserFilter) ([]domain.User, error) {
	return s.users.List(ctx, filter)
}
