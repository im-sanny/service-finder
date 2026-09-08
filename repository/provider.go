package repository

import (
	"database/sql"

	"github.com/im-sanny/service-finder/model"
)

type ProviderRepository interface {
	Create(p *model.Provider) error
	GetAll() ([]model.Provider, error)
	GetById(id int64) (*model.Provider, error)
	Update(p *model.Provider) error
	Patch(id int64, name, phone, location, description *string, service_id *int64) (*model.Provider, error)
	Delete(id int64) error
}

type ppr struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) ProviderRepository {
	return &ppr{db: db}
}

func (h *ppr) Create(p *model.Provider) error {
	err := h.db.QueryRow(`
	INSERT INTO providers (name, phone, location, description, service_id)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID).Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (h *ppr) GetAll() ([]model.Provider, error) {
	rows, err := h.db.Query(`
	SELECT id, name, phone,
	location, description,
	service_id
	FROM providers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var provider []model.Provider
	for rows.Next() {
		var p model.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID); err != nil {
			return nil, err
		}
		provider = append(provider, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (h *ppr) GetById(id int64) (*model.Provider, error) {
	var p model.Provider
	err := h.db.QueryRow(`
	SELECT
	id, name, phone,
	location, description,
	service_id
	FROM providers
	WHERE id=$1`,
		id).Scan(&p.ID, &p.Name,
		&p.Phone, &p.Location,
		&p.Description, &p.ServiceID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (h *ppr) Update(p *model.Provider) error {
	err := h.db.QueryRow(`
	UPDATE providers SET name=$1, phone=$2,
	location=$3, description=$4,
	service_id=$5 WHERE id=$6
	RETURNING id, name, phone,
	location, description, service_id`,
		p.Name, p.Phone, p.Location,
		p.Description, p.ServiceID, p.ID).Scan(&p.ID,
		&p.Name, &p.Phone, &p.Location,
		&p.Description, &p.ServiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	return nil
}

func (h *ppr) Patch(id int64, name, phone, location, description *string, service_id *int64) (*model.Provider, error) {
	var p model.Provider

	err := h.db.QueryRow(`
	UPDATE providers SET
	name=COALESCE($1, name),
	phone=COALESCE($2, phone),
	location=COALESCE($3, location),
	description=COALESCE($4, description),
	service_id=COALESCE($5,	service_id)
	WHERE id=$6
	RETURNING id, name, phone,
	location, description, service_id`,
		p.Name, p.Phone,
		p.Location, p.Description,
		p.ServiceID, id).Scan(&p.ID,
		&p.Name, &p.Phone,
		&p.Location, &p.Description,
		&p.ServiceID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (h *ppr) Delete(id int64) error {
	// Use h.DB.Exec for operations that do not return rows
	res, err := h.db.Exec(`DELETE FROM providers WHERE id=$1`, id)
	if err != nil {
		return err
	}

	// Check how many rows were actually deleted
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}
	return nil
}
