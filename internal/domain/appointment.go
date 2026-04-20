package domain

import "time"

type AppointmentStatus string

const (
	AppointmentPending   AppointmentStatus = "pending"
	AppointmentConfirmed AppointmentStatus = "confirmed"
	AppointmentCancelled AppointmentStatus = "cancelled"
	AppointmentCompleted AppointmentStatus = "completed"
)

type Appointment struct {
	ID                 string            `json:"id"`
	ClientID           string            `json:"client_id"`
	ClientName         string            `json:"client_name"`
	ClientEmail        string            `json:"client_email"`
	ServiceID          string            `json:"service_id"`
	ServiceName        string            `json:"service_name"`
	ProviderID         string            `json:"provider_id"`
	ProviderName       string            `json:"provider_name"`
	StartTime          time.Time         `json:"start_time"`
	EndTime            time.Time         `json:"end_time"`
	Status             AppointmentStatus `json:"status"`
	Notes              string            `json:"notes,omitempty"`
	CancellationReason string            `json:"cancellation_reason,omitempty"`
	ConfirmedAt        *time.Time        `json:"confirmed_at,omitempty"`
	CancelledAt        *time.Time        `json:"cancelled_at,omitempty"`
	CompletedAt        *time.Time        `json:"completed_at,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}
