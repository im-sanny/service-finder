package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/im-sanny/service-finder/model"
)

type ProviderRepository interface {
	Create(p *model.Provider) error
	GetAll() ([]model.Provider, error)
	GetById(id int64) (*model.Provider, error)
	Update(p *model.Provider) error
	Patch(id int64, name, phone, location, description *string, serviceID *int64) (*model.Provider, error)
	Delete(id int64) error
}

type ppr struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) ProviderRepository {
	return &ppr{db: db}
}

func (r *ppr) Create(p *model.Provider) error {
	err := r.db.QueryRow(`
		INSERT INTO providers (name, phone, location, description, service_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	return nil
}

func (r *ppr) GetAll() ([]model.Provider, error) {
	rows, err := r.db.Query(`
		SELECT id, name, phone, location, description, service_id, created_at, updated_at
		FROM providers`)
	if err != nil {
		return nil, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	var providers []model.Provider
	for rows.Next() {
		var p model.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan provider row: %w", err)
		}
		providers = append(providers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during provider iteration: %w", err)
	}

	return providers, nil
}

func (r *ppr) GetById(id int64) (*model.Provider, error) {
	var p model.Provider
	err := r.db.QueryRow(`
		SELECT id, name, phone, location, description, service_id, created_at, updated_at
		FROM providers WHERE id=$1`, id,
	).Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %d: %w", id, err)
	}
	return &p, nil
}

func (r *ppr) Update(p *model.Provider) error {
	err := r.db.QueryRow(`
		UPDATE providers
		SET name=$1, phone=$2, location=$3, description=$4, service_id=$5, created_at=$6, updated_at=$7
		WHERE id=$8
		RETURNING id, name, phone, location, description, service_id, created_at, updated_at`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID, p.CreatedAt, p.UpdatedAt, p.ID,
	).Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to update provider %d: %w", p.ID, err)
	}
	return nil
}

func (r *ppr) Patch(id int64, name, phone, location, description *string, serviceID *int64) (*model.Provider, error) {
	var p model.Provider

	err := r.db.QueryRow(`
		UPDATE providers SET
			name=COALESCE($1, name),
			phone=COALESCE($2, phone),
			location=COALESCE($3, location),
			description=COALESCE($4, description),
			service_id=COALESCE($5, service_id)
		WHERE id=$6
		RETURNING id, name, phone, location, description, service_id, created_at, updated_at`,
		name, phone, location, description, serviceID, id,
	).Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to patch provider %d: %w", id, err)
	}
	return &p, nil
}

func (r *ppr) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM providers WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete for provider %d: %w", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected for provider %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
