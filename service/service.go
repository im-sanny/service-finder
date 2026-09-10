package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/repository"
)

// Domain errors shared by all business logic.
// Handlers use errors.Is() to map these to HTTP status codes.
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
)

// Service defines the business logic for managing Service entities.
// Handlers depend on THIS interface, never on the repository directly.
type Service interface {
	Create(s *model.Service) error
	GetAll() ([]model.Service, error)
	GetByID(id int64) (*model.Service, error)
	Update(s *model.Service) error
	Patch(id int64, name, description *string) (*model.Service, error)
	Delete(id int64) error
}

// service is the PRIVATE implementation.
// Outside code can only access it through the Service interface.
type service struct {
	repo repository.ServiceRepository
}

// NewService wires the business logic to its data source.
func NewService(repo repository.ServiceRepository) Service {
	return &service{repo: repo}
}

func (s *service) Create(svc *model.Service) error {
	// Business rule: name is required.
	// model.Service uses pointers, so we check both nil and empty string.
	if svc.Name == nil || *svc.Name == "" {
		return fmt.Errorf("service name required: %w", ErrInvalidInput)
	}
	if svc.Description == nil || *svc.Description == "" {
		return fmt.Errorf("service description required: %w", ErrInvalidInput)
	}
	if err := s.repo.Create(svc); err != nil {
		return fmt.Errorf("create service: %w", err)
	}
	return nil
}

func (s *service) GetAll() ([]model.Service, error) {
	services, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get all services: %w", err)
	}
	return services, nil
}

func (s *service) GetByID(id int64) (*model.Service, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid service id %d: %w", id, ErrInvalidInput)
	}

	svc, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("service %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("get service %d: %w", id, err)
	}
	return svc, nil
}

func (s *service) Update(svc *model.Service) error {
	if svc.ID <= 0 {
		return fmt.Errorf("invalid service id %d: %w", svc.ID, ErrInvalidInput)
	}
	if svc.Name == nil || *svc.Name == "" {
		return fmt.Errorf("service name required: %w", ErrInvalidInput)
	}
	if svc.Description == nil || *svc.Description == "" {
		return fmt.Errorf("service description required: %w", ErrInvalidInput)
	}

	if err := s.repo.Update(svc); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service %d: %w", svc.ID, ErrNotFound)
		}
		return fmt.Errorf("update service %d: %w", svc.ID, err)
	}
	return nil
}

func (s *service) Patch(id int64, name, description *string) (*model.Service, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid service id %d: %w", id, ErrInvalidInput)
	}

	svc, err := s.repo.Patch(id, name, description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("service %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("patch service %d: %w", id, err)
	}
	return svc, nil
}

func (s *service) Delete(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid service id %d: %w", id, ErrInvalidInput)
	}
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service %d: %w", id, ErrNotFound)
		}
		return fmt.Errorf("delete service %d: %w", id, err)
	}
	return nil
}
