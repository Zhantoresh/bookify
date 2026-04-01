package domain

import (
	"errors"
	"time"
)

type Role string

const (
	RoleClient     Role = "client"
	RoleSpecialist Role = "specialist"
	RoleAdmin      Role = "admin"
)

// Validate проверяет, что роль существует
func (r Role) Validate() error {
	switch r {
	case RoleClient, RoleSpecialist, RoleAdmin:
		return nil
	default:
		return errors.New("неверная роль: " + string(r))
	}
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// IsAdmin — быстрый способ проверить права
func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
