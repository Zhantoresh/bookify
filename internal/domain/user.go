package domain

import "time"

type Role string

const (
	RoleClient     Role = "client"
	RoleSpecialist Role = "specialist"
	RoleAdmin      Role = "admin"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
