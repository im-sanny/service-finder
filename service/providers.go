package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/repository"
)

type Providers interface {
	Create(p *model.Provider) error
	GetAll() ([]model.Provider, error)
	GetByID(id int64) (*model.Provider, error)
	Update(p *model.Provider) error
	Patch(id int64, name, phone, location, description *string, serviceID *int64) (*model.Provider, error)
	Delete(id int64) error

	CreateBatch(p []*model.Provider) ([]*model.Provider, error)
	DeleteBatch(ids []int64) (int64, error)
}

type provider struct {
	repo        repository.ProviderRepository
	serviceRepo repository.ServiceRepository
}

func NewProvider(repo repository.ProviderRepository, serviceRepo repository.ServiceRepository) Providers {
	return &provider{
		repo:        repo,
		serviceRepo: serviceRepo,
	}
}

func (s *provider) DeleteBatch(ids []int64) (int64, error) {
	// 1. Validate: slice not empty
	if len(ids) == 0 {
		return 0, fmt.Errorf("ids list is empty: %w", ErrInvalidInput)
	}

	// 2. Validate all ids are positive
	for i, id := range ids {
		if id <= 0 {
			return 0, fmt.Errorf("id at index %d is invalid: %w", i, ErrInvalidInput)
		}
	}

	// 3. Start transaction
	tx, err := s.repo.BeginTx()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// 4. Batch delete
	deleted, err := s.repo.DeleteBatch(tx, ids)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed: %w", err)
	}

	// 5. Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return deleted, nil
}

// CreateBatch creates multiple providers atomically
func (s *provider) CreateBatch(prov []*model.Provider) ([]*model.Provider, error) {
	// 1. Validate: slice is not empty
	if len(prov) == 0 {
		return nil, fmt.Errorf("provider list is empty: %w", ErrInvalidInput)
	}

	// 2. Validate EACH provider before starting transaction
	for i, v := range prov {
		if v.Name == nil || *v.Name == "" {
			return nil, fmt.Errorf("provider at index %d: name is required: %w", i, ErrInvalidInput)
		}
	}

	// 3. Verify all ServiceIDs exist (cross-entity validation)
	for i, p := range prov {
		if p.ServiceID != nil {
			_, err := s.serviceRepo.GetById(*p.ServiceID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, fmt.Errorf("provider at index %d: service %d does not exist: %w",
						i, *p.ServiceID, ErrInvalidInput)
				}
				return nil, fmt.Errorf("provider at index %d: failed to verify service %d: %w",
					i, *p.ServiceID, err)
			}
		}
	}

	// 4. Start transaction
	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback on any error (safe to call even after commit)h
	defer tx.Rollback()

	// 5. Batch insert
	if err := s.repo.CreateBatch(tx, prov); err != nil {
		return nil, fmt.Errorf("batch create failed: %w", err)
	}

	// 6. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return prov, nil
}

func (s *provider) Create(p *model.Provider) error {
	if p.Name == nil || *p.Name == "" {
		return fmt.Errorf("provider name required: %w", ErrInvalidInput)
	}
	if p.ServiceID != nil {
		_, err := s.serviceRepo.GetById(*p.ServiceID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("service %d does not exists: %w", *p.ServiceID, ErrNotFound)
			}
			return fmt.Errorf("failed to verify provider %d: %w", *p.ServiceID, err)
		}
	}

	if err := s.repo.Create(p); err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	return nil
}

func (s *provider) GetAll() ([]model.Provider, error) {
	providers, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all providers: %w", err)
	}
	return providers, nil
}

func (s *provider) GetByID(id int64) (*model.Provider, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid provider id %d: %w", id, ErrInvalidInput)
	}

	prov, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("provider %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("get provider %d: %w", id, err)
	}
	return prov, nil
}

func (s *provider) Update(p *model.Provider) error {
	if p.ID <= 0 {
		return fmt.Errorf("invalid provider id %d: %w", p.ID, ErrInvalidInput)
	}
	if p.Name == nil || *p.Name == "" {
		return fmt.Errorf("provider name is required: %w", ErrInvalidInput)
	}
	if p.ServiceID == nil {
		return fmt.Errorf("service id is required: %w", ErrInvalidInput)
	}

	_, err := s.serviceRepo.GetById(*p.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service %d does not exist: %w", *p.ServiceID, ErrInvalidInput)
		}
		return fmt.Errorf("failed to verify service %d: %w", *p.ServiceID, err)
	}

	if err := s.repo.Update(p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("provider %d: %w", p.ID, ErrNotFound)
		}
		return fmt.Errorf("failed to update provider %d: %w", p.ID, err)
	}
	return nil
}

func (s *provider) Patch(id int64, name, phone, location, description *string, serviceID *int64) (*model.Provider, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid provider id %d: %w", id, ErrInvalidInput)
	}

	if serviceID != nil {
		_, err := s.serviceRepo.GetById(*serviceID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("service %d does not exist: %w", *serviceID, ErrInvalidInput)
			}
			return nil, fmt.Errorf("failed to verify service %d: %w", *serviceID, err)
		}
	}

	pvr, err := s.repo.Patch(id, name, phone, location, description, serviceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("provider %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to patch provider %d: %w", id, err)
	}
	return pvr, nil
}

func (s *provider) Delete(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid provider id %d: %w", id, ErrInvalidInput)
	}

	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("provider %d: %w", id, ErrNotFound)
		}
		return fmt.Errorf("failed to delete provider %d: %w", id, err)
	}
	return nil
}
