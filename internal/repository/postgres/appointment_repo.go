package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *AppointmentRepository) CheckOverlap(ctx context.Context, tx *sql.Tx, serviceID string, start, end time.Time) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1 FROM appointments
			WHERE service_id = $1
			  AND status IN ('pending', 'confirmed')
			  AND start_time < $2
			  AND end_time > $3
		)`
	var exists bool
	if err := tx.QueryRowContext(ctx, query, serviceID, end, start).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AppointmentRepository) Create(ctx context.Context, tx *sql.Tx, appointment *domain.Appointment) error {
	const query = `
		INSERT INTO appointments (client_id, service_id, start_time, end_time, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	return tx.QueryRowContext(ctx, query,
		appointment.ClientID, appointment.ServiceID, appointment.StartTime, appointment.EndTime, appointment.Status, nullableString(appointment.Notes),
	).Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)
}

func (r *AppointmentRepository) List(ctx context.Context, filter repository.AppointmentFilter) ([]domain.Appointment, repository.Pagination, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	var conditions []string
	var args []any
	argPos := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("a.status = $%d", argPos))
		args = append(args, filter.Status)
		argPos++
	}
	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("a.start_time >= $%d", argPos))
		args = append(args, *filter.FromDate)
		argPos++
	}
	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("a.start_time <= $%d", argPos))
		args = append(args, *filter.ToDate)
		argPos++
	}
	if filter.ProviderID != "" {
		conditions = append(conditions, fmt.Sprintf("s.provider_id = $%d", argPos))
		args = append(args, filter.ProviderID)
		argPos++
	}
	if filter.ClientID != "" {
		conditions = append(conditions, fmt.Sprintf("a.client_id = $%d", argPos))
		args = append(args, filter.ClientID)
		argPos++
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := `
		SELECT COUNT(*)
		FROM appointments a
		JOIN services s ON s.id = a.service_id` + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, repository.Pagination{}, err
	}

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, filter.Limit, offset)
	query := fmt.Sprintf(`
		SELECT a.id, a.client_id, cu.full_name, cu.email, a.service_id, s.name, s.provider_id, pu.full_name,
		       a.start_time, a.end_time, a.status, COALESCE(a.notes, ''), COALESCE(a.cancellation_reason, ''),
		       a.confirmed_at, a.cancelled_at, a.completed_at, a.created_at, a.updated_at
		FROM appointments a
		JOIN users cu ON cu.id = a.client_id
		JOIN services s ON s.id = a.service_id
		JOIN users pu ON pu.id = s.provider_id
		%s
		ORDER BY a.start_time DESC
		LIMIT $%d OFFSET $%d`, where, argPos, argPos+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, repository.Pagination{}, err
	}
	defer rows.Close()
	items, err := scanAppointments(rows)
	return items, buildPagination(filter.Page, filter.Limit, total), err
}

func (r *AppointmentRepository) ListByActor(ctx context.Context, actorID string, role domain.Role, page, limit int) ([]domain.Appointment, repository.Pagination, error) {
	filter := repository.AppointmentFilter{Page: page, Limit: limit}
	if role == domain.RoleClient {
		filter.ClientID = actorID
	} else {
		filter.ProviderID = actorID
	}
	return r.List(ctx, filter)
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	query := `
		SELECT a.id, a.client_id, cu.full_name, cu.email, a.service_id, s.name, s.provider_id, pu.full_name,
		       a.start_time, a.end_time, a.status, COALESCE(a.notes, ''), COALESCE(a.cancellation_reason, ''),
		       a.confirmed_at, a.cancelled_at, a.completed_at, a.created_at, a.updated_at
		FROM appointments a
		JOIN users cu ON cu.id = a.client_id
		JOIN services s ON s.id = a.service_id
		JOIN users pu ON pu.id = s.provider_id
		WHERE a.id = $1`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items, err := scanAppointments(rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, domain.ErrNotFound
	}
	return &items[0], nil
}

func (r *AppointmentRepository) UpdateStatus(ctx context.Context, id string, status domain.AppointmentStatus, reason string, changedAt time.Time) error {
	var (
		query string
		args  []any
	)

	switch status {
	case domain.AppointmentConfirmed:
		query = `
			UPDATE appointments
			SET status = $2,
			    confirmed_at = $3,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1`
		args = []any{id, string(status), changedAt}
	case domain.AppointmentCancelled:
		query = `
			UPDATE appointments
			SET status = $2,
			    cancellation_reason = $3,
			    cancelled_at = $4,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1`
		args = []any{id, string(status), nullableString(reason), changedAt}
	case domain.AppointmentCompleted:
		query = `
			UPDATE appointments
			SET status = $2,
			    completed_at = $3,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1`
		args = []any{id, string(status), changedAt}
	default:
		query = `
			UPDATE appointments
			SET status = $2,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1`
		args = []any{id, string(status)}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *AppointmentRepository) GetAppointmentsByDateRange(ctx context.Context, start, end time.Time) ([]domain.Appointment, error) {
	query := `
		SELECT a.id, a.client_id, cu.full_name, cu.email, a.service_id, s.name, s.provider_id, pu.full_name,
		       a.start_time, a.end_time, a.status, COALESCE(a.notes, ''), COALESCE(a.cancellation_reason, ''),
		       a.confirmed_at, a.cancelled_at, a.completed_at, a.created_at, a.updated_at
		FROM appointments a
		JOIN users cu ON cu.id = a.client_id
		JOIN services s ON s.id = a.service_id
		JOIN users pu ON pu.id = s.provider_id
		WHERE a.start_time >= $1 AND a.start_time < $2`
	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAppointments(rows)
}

func scanAppointments(rows *sql.Rows) ([]domain.Appointment, error) {
	var items []domain.Appointment
	for rows.Next() {
		var item domain.Appointment
		if err := rows.Scan(
			&item.ID, &item.ClientID, &item.ClientName, &item.ClientEmail, &item.ServiceID, &item.ServiceName, &item.ProviderID,
			&item.ProviderName, &item.StartTime, &item.EndTime, &item.Status, &item.Notes, &item.CancellationReason,
			&item.ConfirmedAt, &item.CancelledAt, &item.CompletedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func buildPagination(page, limit, total int) repository.Pagination {
	pages := 0
	if total > 0 {
		pages = (total + limit - 1) / limit
	}
	return repository.Pagination{Page: page, Limit: limit, Total: total, Pages: pages}
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
