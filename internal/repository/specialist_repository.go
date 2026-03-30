package repository

import (
	"database/sql"

	"github.com/bookify/internal/domain"
)

type SpecialistRepository struct {
	db *sql.DB
}

func NewSpecialistRepository(db *sql.DB) *SpecialistRepository {
	return &SpecialistRepository{db: db}
}

// GetAll returns all specialists (users with role 'specialist')
func (r *SpecialistRepository) GetAll() ([]domain.Specialist, error) {
	query := `SELECT id, name FROM users WHERE role = $1 ORDER BY id`

	rows, err := r.db.Query(query, domain.RoleSpecialist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialists []domain.Specialist
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		specialists = append(specialists, domain.Specialist{
			ID:   id,
			Name: name,
			Type: "specialist", // Default type
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialists, nil
}

// GetByID returns a specialist by user ID
func (r *SpecialistRepository) GetByID(id int) (*domain.Specialist, error) {
	query := `SELECT id, name FROM users WHERE id = $1 AND role = $2`

	var specID int
	var name string
	err := r.db.QueryRow(query, id, domain.RoleSpecialist).Scan(&specID, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Specialist{
		ID:   specID,
		Name: name,
		Type: "specialist",
	}, nil
}
