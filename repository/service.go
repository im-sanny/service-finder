package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/im-sanny/service-finder/model"
)

// ServiceRepository defines the contract for service data operations.
type ServiceRepository interface {
	GetAll() ([]model.Service, error)
	GetById(id int64) (*model.Service, error)
	Create(s *model.Service) error
	Update(s *model.Service) error
	Patch(id int64, name, description *string) (*model.Service, error)
	Delete(id int64) error
}

type psr struct {
	db *sql.DB // use pointer to avoid copying and share single, safe connection pool
}

// NewServiceRepository creates a new instance, reusing the single, safe connection pool.
func NewServiceRepository(db *sql.DB) ServiceRepository {
	return &psr{db: db}
}

func (r *psr) Create(s *model.Service) error {
	err := r.db.QueryRow(`
		INSERT INTO services (name, description)
		VALUES ($1, $2) RETURNING id`, // RETURNING id gives you one newly-created ID.
		s.Name, s.Description).Scan(&s.ID) // The & means you're giving Scan the memory addresses where it should put the values.
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	return nil
}

func (r *psr) GetAll() ([]model.Service, error) {
	rows, err := r.db.Query(`SELECT id, name, description, created_at, updated_at FROM services;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all services: %w", err)
	}

	defer rows.Close() // When this handler finishes, close the rows automatically.

	var services []model.Service
	// The loop basically means:
	// "Give me the first row → scan it → put it in my slice.
	// Give me the next row → scan it → put it in my slice.
	// Keep going until there are no more rows."
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan service row: %w", err)
		}
		services = append(services, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return services, nil
}

func (r *psr) GetById(id int64) (*model.Service, error) {
	var s model.Service
	err := r.db.QueryRow(`
		SELECT id, name, description, created_at, updated_at
		FROM services WHERE id=$1`, id).Scan(
		&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		// FIXED: Returning a pointer to an empty struct on error is incorrect.
		// Return nil so the caller can handle the error (like sql.ErrNoRows) properly.
		return nil, fmt.Errorf("failed to get service %d: %w", id, err)
	}

	return &s, nil
}

func (r *psr) Update(s *model.Service) error {
	err := r.db.QueryRow(`
		UPDATE services
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at`,
		s.Name, s.Description, s.ID, // The SQL says what values it needs. Go supplies them.
	).Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		// Because we use RETURNING, 0 rows affected results in sql.ErrNoRows from Scan().
		// We handle it explicitly to satisfy the "RowsAffected" requirement.
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to update service %d: %w", s.ID, err)
	}

	return nil
}

func (r *psr) Patch(id int64, name, description *string) (*model.Service, error) {
	var s model.Service

	err := r.db.QueryRow(`
		UPDATE services
		SET name=COALESCE($1, name), description=COALESCE($2, description)
		WHERE id=$3
		RETURNING id, name, description, created_at, updated_at`,
		name, description, id).Scan(
		&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to patch service %d: %w", id, err)
	}
	return &s, nil
}

func (r *psr) Delete(id int64) error {
	// 1. Use Exec instead of QueryRow for operations that don't need to scan returning data
	res, err := r.db.Exec(`DELETE FROM services WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete for service %d: %w", id, err)
	}

	// 2. Check how many rows were actually affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected for service %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
