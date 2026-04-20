package domain

import "time"

type Service struct {
	ID              string     `json:"id"`
	ProviderID      string     `json:"provider_id"`
	ProviderName    string     `json:"provider_name,omitempty"`
	ProviderEmail   string     `json:"provider_email,omitempty"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	Price           float64    `json:"price"`
	DurationMinutes int        `json:"duration_minutes"`
	IsActive        bool       `json:"is_active"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
