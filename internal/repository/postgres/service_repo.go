package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	const query = `
		INSERT INTO services (provider_id, name, description, price, duration_minutes, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, service.ProviderID, service.Name, nullableString(service.Description), service.Price, service.DurationMinutes, service.IsActive).
		Scan(&service.ID, &service.CreatedAt, &service.UpdatedAt)
}

func (r *ServiceRepository) List(ctx context.Context, filter repository.ServiceFilter) ([]domain.Service, repository.Pagination, error) {
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

	if filter.ProviderID != "" {
		conditions = append(conditions, fmt.Sprintf("s.provider_id = $%d", argPos))
		args = append(args, filter.ProviderID)
		argPos++
	}
	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("s.price >= $%d", argPos))
		args = append(args, *filter.MinPrice)
		argPos++
	}
	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("s.price <= $%d", argPos))
		args = append(args, *filter.MaxPrice)
		argPos++
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("s.name ILIKE $%d", argPos))
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}
	if filter.OnlyActive {
		conditions = append(conditions, "s.is_active = true AND s.deleted_at IS NULL")
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := `SELECT COUNT(*) FROM services s` + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, repository.Pagination{}, err
	}

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, filter.Limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT s.id, s.provider_id, u.full_name, u.email, s.name, COALESCE(s.description, ''), s.price,
		       s.duration_minutes, s.is_active, s.deleted_at, s.created_at, s.updated_at
		FROM services s
		JOIN users u ON u.id = s.provider_id
		%s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d`, where, argPos, argPos+1)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, repository.Pagination{}, err
	}
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var item domain.Service
		if err := rows.Scan(
			&item.ID, &item.ProviderID, &item.ProviderName, &item.ProviderEmail, &item.Name, &item.Description,
			&item.Price, &item.DurationMinutes, &item.IsActive, &item.DeletedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, repository.Pagination{}, err
		}
		services = append(services, item)
	}

	return services, buildPagination(filter.Page, filter.Limit, total), rows.Err()
}

func (r *ServiceRepository) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	const query = `
		SELECT s.id, s.provider_id, u.full_name, u.email, s.name, COALESCE(s.description, ''), s.price,
		       s.duration_minutes, s.is_active, s.deleted_at, s.created_at, s.updated_at
		FROM services s
		JOIN users u ON u.id = s.provider_id
		WHERE s.id = $1`

	var service domain.Service
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&service.ID, &service.ProviderID, &service.ProviderName, &service.ProviderEmail, &service.Name, &service.Description,
		&service.Price, &service.DurationMinutes, &service.IsActive, &service.DeletedAt, &service.CreatedAt, &service.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &service, nil
}

func (r *ServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	const query = `
		UPDATE services
		SET name = $2, description = $3, price = $4, duration_minutes = $5,
		    is_active = $6, deleted_at = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at`
	if err := r.db.QueryRowContext(ctx, query,
		service.ID, service.Name, nullableString(service.Description), service.Price, service.DurationMinutes,
		service.IsActive, service.DeletedAt,
	).Scan(&service.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *ServiceRepository) HasFutureAppointments(ctx context.Context, serviceID string, now time.Time) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1 FROM appointments
			WHERE service_id = $1
			  AND start_time > $2
			  AND status IN ('pending', 'confirmed')
		)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, serviceID, now).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
