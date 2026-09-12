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

func (s *provider) Create(p *model.Provider) error {
	if p.Name == nil || *p.Name == "" {
		return fmt.Errorf("service name required: %w", ErrInvalidInput)
	}
	if p.ServiceID != nil {
		_, err := s.serviceRepo.GetById(*p.ServiceID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("service id %d does not exists: %w", *p.ServiceID, ErrNotFound)
			}
			return fmt.Errorf("failed to verify service %d: %w", *p.ServiceID, err)
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
