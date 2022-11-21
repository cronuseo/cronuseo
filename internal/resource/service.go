package resource

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Resource, error)
	Query(ctx context.Context, org_id string) ([]Resource, error)
	Create(ctx context.Context, org_id string, input CreateResourceRequest) (Resource, error)
	Update(ctx context.Context, org_id string, id string, input UpdateResourceRequest) (Resource, error)
	Delete(ctx context.Context, org_id string, id string) (Resource, error)
}

type Resource struct {
	entity.Resource
}

type CreateResourceRequest struct {
	Key   string `json:"resource_key" db:"resource_key"`
	Name  string `json:"name" db:"name"`
	OrgID string `json:"-" db:"org_id"`
}

func (m CreateResourceRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Key, validation.Required),
	)
}

type UpdateResourceRequest struct {
	Name string `json:"name" db:"name"`
}

func (m UpdateResourceRequest) Validate() error {
	return validation.ValidateStruct(&m)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) Get(ctx context.Context, org_id string, id string) (Resource, error) {
	resource, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		return Resource{}, err
	}
	return Resource{resource}, nil
}

func (s service) Create(ctx context.Context, org_id string, req CreateResourceRequest) (Resource, error) {
	if err := req.Validate(); err != nil {
		return Resource{}, err
	}
	id := entity.GenerateID()
	err := s.repo.Create(ctx, org_id, entity.Resource{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		return Resource{}, err
	}
	return s.Get(ctx, org_id, id)
}

func (s service) Update(ctx context.Context, org_id string, id string, req UpdateResourceRequest) (Resource, error) {
	if err := req.Validate(); err != nil {
		return Resource{}, err
	}

	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		return resource, err
	}
	resource.Name = req.Name
	if err := s.repo.Update(ctx, org_id, resource.Resource); err != nil {
		return resource, err
	}
	return resource, nil
}

func (s service) Delete(ctx context.Context, org_id string, id string) (Resource, error) {
	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		return Resource{}, err
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		return Resource{}, err
	}
	return resource, nil
}

func (s service) Query(ctx context.Context, org_id string) ([]Resource, error) {
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		return nil, err
	}
	result := []Resource{}
	for _, item := range items {
		result = append(result, Resource{item})
	}
	return result, nil
}
