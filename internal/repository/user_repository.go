package repository

import (
	"database/sql"

	"github.com/bookify/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (email, name, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(query, user.Email, user.Name, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, name, password_hash, role, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, name, password_hash, role, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetAllSpecialists returns all users with role 'specialist'
func (r *UserRepository) GetAllSpecialists() ([]domain.User, error) {
	query := `SELECT id, email, name, password_hash, role, created_at FROM users WHERE role = $1 ORDER BY id`

	rows, err := r.db.Query(query, domain.RoleSpecialist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialists []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		specialists = append(specialists, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialists, nil
}
