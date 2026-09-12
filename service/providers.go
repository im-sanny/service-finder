package service

import (
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

