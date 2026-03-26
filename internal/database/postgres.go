package database

import (
    "bookify/internal/domain"
    "database/sql"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}


func (r *UserRepository) CreateUser(user *domain.User) error {
    query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, created_at`
    return r.db.QueryRow(query, user.Email, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt)
}


func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
    user := &domain.User{}
    query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return user, nil
}