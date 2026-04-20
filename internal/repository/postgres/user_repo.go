package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bookify/internal/domain"
	"github.com/bookify/internal/repository"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (email, password_hash, full_name, role, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	if err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.FullName, user.Role, nullableString(user.Phone)).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, full_name, role, COALESCE(phone, ''), created_at, updated_at
		FROM users WHERE id = $1`
	return scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, full_name, role, COALESCE(phone, ''), created_at, updated_at
		FROM users WHERE email = $1`
	return scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *UserRepository) List(ctx context.Context, filter repository.UserFilter) ([]domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, COALESCE(phone, ''), created_at, updated_at
		FROM users`
	args := []any{}
	if filter.Role != "" {
		query += ` WHERE role = $1`
		args = append(args, filter.Role)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.Phone, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func scanUser(scanner interface{ Scan(dest ...any) error }) (*domain.User, error) {
	var user domain.User
	if err := scanner.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.Phone, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
