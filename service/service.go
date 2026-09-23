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
	GetAll(page, limit int) ([]model.Service, error)
	GetByID(id int64) (*model.Service, error)
	Update(s *model.Service) error
	Patch(id int64, name, description *string) (*model.Service, error)
	Delete(id int64) error

	CreateBatch(s []*model.Service) ([]*model.Service, error)
	DeleteBatch(ids []int64) (int64, error)
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

func (s *service) CreateBatch(svc []*model.Service) ([]*model.Service, error) {
	if len(svc) == 0 {
		return nil, fmt.Errorf("service list is empty: %w", ErrInvalidInput)
	}

	for i, v := range svc {
		if v.Name == nil || *v.Name == "" {
			return nil, fmt.Errorf("service at index %d: name is required: %w", i, ErrInvalidInput)
		}
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer a cleanup function that checks the named 'err' variable
	defer func() {
		if err != nil {
			// If we are exiting with an error, try to rollback to clean up DB locks
			tx.Rollback()
		}
	}()

	if err := s.repo.CreateBatch(tx, svc); err != nil {
		return nil, fmt.Errorf("batch insert failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return svc, nil
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

func (s *service) GetAll(page, limit int) ([]model.Service, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	services, err := s.repo.GetAll(page, limit)
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

func (s *service) DeleteBatch(ids []int64) (int64, error) {
	// 1. Validate: slice not empty
	if len(ids) == 0 {
		return 0, fmt.Errorf("ids list is empty: %w", ErrInvalidInput)
	}

	// 2. Validate all IDs are Positive
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

	// 4. Defer a cleanup function that checks the named 'err' variable
	defer func() {
		if err != nil {
			// If we are exiting with an error, try to rollback to clean up DB locks
			tx.Rollback()
		}
	}()

	// 5. Batch delete
	deleted, err := s.repo.DeleteBatch(tx, ids)
	if err != nil {
		return 0, fmt.Errorf("batch delete failed: %w", err)
	}

	// 6. Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return deleted, nil
}
