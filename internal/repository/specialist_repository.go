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

func (r *SpecialistRepository) GetAll() ([]domain.Specialist, error) {
	query := `SELECT id, name, type FROM specialists ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specialists []domain.Specialist
	for rows.Next() {
		var spec domain.Specialist
		err := rows.Scan(&spec.ID, &spec.Name, &spec.Type)
		if err != nil {
			return nil, err
		}
		specialists = append(specialists, spec)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return specialists, nil
}

func (r *SpecialistRepository) GetByID(id int) (*domain.Specialist, error) {
	query := `SELECT id, name, type FROM specialists WHERE id = $1`

	var spec domain.Specialist
	err := r.db.QueryRow(query, id).Scan(&spec.ID, &spec.Name, &spec.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &spec, nil
}
