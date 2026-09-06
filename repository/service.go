package repository

import (
	"database/sql"
	"fmt"

	"github.com/im-sanny/service-finder/model"
)

// ServiceRepository defines the contract for service data operations.
type ServiceRepository interface {
	GetAll() ([]model.Service, error)
	GetById(id int64) (*model.Service, error)
	Create(s *model.Service) error
	Update(s *model.Service) error
	Patch(id int64, name, description string) (*model.Service, error)
	Delete(id int64) error
}

type psr struct {
	db *sql.DB // use pointer to avoid copying and share single, safe connection pool
}

// NewServiceRepository creates a new instance, reusing the single, safe connection pool.
func NewServiceRepository(db *sql.DB) ServiceRepository {
	return &psr{db: db}
}

func (h *psr) Create(s *model.Service) error {
	err := h.db.QueryRow(`
		INSERT INTO
		services (name, description)
		VALUES ($1, $2) RETURNING id`, // RETURNING id gives you one newly-created ID.
		s.Name, s.Description).Scan(&s.ID) // The & means you're giving Scan the memory addresses where it should put the values.
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	return nil
}

func (h *psr) GetAll() ([]model.Service, error) {
	rows, err := h.db.Query(`SELECT id, name, description, created_at, updated_at FROM services;`)
	if err != nil {
		return nil, err
	}

	defer rows.Close() // When this handler finishes, close the rows automatically.

	services := make([]model.Service, 0)
	// The loop basically means:
	// "Give me the first row → scan it → put it in my slice.
	// Give me the next row → scan it → put it in my slice.
	// Keep going until there are no more rows."
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (h *psr) GetById(id int64) (*model.Service, error) {
	var s model.Service
	err := h.db.QueryRow(`
	SELECT
	id, name, description,
	created_at, updated_at
	FROM services WHERE id=$1`, id).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return &model.Service{}, err
	}

	return &s, nil
}

func (h *psr) Update(s *model.Service) error {
	err := h.db.QueryRow(`
	UPDATE services
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description,
		created_at, updated_at`,
		s.Name, s.Description, s.ID, // The SQL says what values it needs. Go supplies them.
	).Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update service %d: %w", s.ID, err)
	}

	return nil
}

func (h *psr) Patch(id int64, name, description string) (*model.Service, error) {
	var s model.Service

	err := h.db.QueryRow(`
	UPDATE services
	SET name=COALESCE($1, name),
	description=COALESCE($2, description)
	WHERE id=$3
	RETURNING
	id, name, description,
	created_at, updated_at`,
		name, description, id).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (h *psr) Delete(id int64) error {

	var deletedId int
	err := h.db.QueryRow(`DELETE FROM services WHERE id=$1 RETURNING id`, id).Scan(&deletedId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("failed to delete service %d: %w", id, err)
		}
	}
	return nil
}
