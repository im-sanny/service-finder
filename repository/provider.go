package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/im-sanny/service-finder/model"
	"github.com/lib/pq"
)

type ProviderRepository interface {
	Create(p *model.Provider) error
	GetAll(page, limit int, filters map[string]string) ([]model.Provider, int64, error)
	GetById(id int64) (*model.Provider, error)
	Update(p *model.Provider) error
	Patch(id int64, name, phone, location, description *string, serviceID *int64) (*model.Provider, error)
	Delete(id int64) error

	// Batch operations with transaction support
	BeginTx() (*sql.Tx, error)
	CreateBatch(tx *sql.Tx, p []*model.Provider) error
	DeleteBatch(tx *sql.Tx, ids []int64) (int64, error)
}

type ppr struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) ProviderRepository {
	return &ppr{db: db}
}

func (r *ppr) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

// CreateBatch inserts multiple providers using a prepared statement within a transaction
func (r *ppr) CreateBatch(tx *sql.Tx, providers []*model.Provider) error {
	// Prepare once, execute many times
	stmt, err := tx.Prepare(`
	INSERT INTO providers (name, phone, location, description, service_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`) // // Prepare the statement ONCE
	if err != nil {
		return fmt.Errorf("failed to prepare batch insert: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement for each provider
	for _, p := range providers {
		err := stmt.QueryRow(
			p.Name, p.Phone, p.Location, p.Description, p.ServiceID,
		).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert provider: %w", err)
		}
	}
	return nil
}

func (r *ppr) Create(p *model.Provider) error {
	err := r.db.QueryRow(`
		INSERT INTO providers (name, phone, location, description, service_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	return nil
}

func (r *ppr) GetAll(page, limit int, filters map[string]string) ([]model.Provider, int64, error) {
	query := `SELECT id, name, phone, location, description, service_id, created_at, updated_at FROM providers`
	countQuery := `SELECT COUNT(*) FROM providers`

	var args []interface{}
	argIndex := 1
	whereClauses := []string{}

	if loc, ok := filters["location"]; ok && loc != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("location ILIKE $%d", argIndex))
		args = append(args, "%"+loc+"%")
		argIndex++
	}

	if sid, ok := filters["service_id"]; ok && sid != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("service_id = $%d", argIndex))
		args = append(args, sid)
		argIndex++
	}

	if len(whereClauses) > 0 {
		whereSQL := " WHERE " + strings.Join(whereClauses, " AND ")
		query += whereSQL
		countQuery += whereSQL
	}

	query += fmt.Sprintf("ORDER BY id ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, (page-1)*limit)

	var total int64
	err := r.db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count providers: %w", err)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query providers: %w", err)
	}
	defer rows.Close()

	var providers []model.Provider
	for rows.Next() {
		var p model.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan provider row: %w", err)
		}
		providers = append(providers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error during provider iteration: %w", err)
	}

	return providers, total, nil
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
		SET
			name = $1,
			phone = $2,
			location = $3,
			description = $4,
			service_id = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, name, phone, location, description, service_id, created_at, updated_at`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID, p.ID,
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
			name = COALESCE($1, name),
			phone = COALESCE($2, phone),
			location = COALESCE($3, location),
			description = COALESCE($4, description),
			service_id = COALESCE($5, service_id),
			updated_at = NOW()
		WHERE id = $6
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

// DeleteBatch deletes multiple providers in a single query
func (r *ppr) DeleteBatch(tx *sql.Tx, ids []int64) (int64, error) {
	// PostgreSQL's ANY() operator with pq.Array for efficient batch delete
	res, err := tx.Exec(`DELETE FROM providers WHERE id = ANY($1)`, pq.Array(ids))
	if err != nil {
		return 0, fmt.Errorf("failed to batch delete: %w", err)
	}
	return res.RowsAffected()
}
