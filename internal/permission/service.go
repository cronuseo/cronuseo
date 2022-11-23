package permission

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Permission, error)
	Query(ctx context.Context, org_id string) ([]Permission, error)
	Create(ctx context.Context, org_id string, input CreateResourceRequest) (Permission, error)
	Update(ctx context.Context, org_id string, id string, input UpdateResourceRequest) (Permission, error)
	Delete(ctx context.Context, org_id string, id string) (Permission, error)
}

type Permission struct {
	entity.Permission
}

type CreateResourceRequest struct {
	Key        string `json:"permission_key" db:"permission_key"`
	Name       string `json:"name" db:"name"`
	ResourceID string `json:"-" db:"resource_id"`
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

func (s service) Get(ctx context.Context, org_id string, id string) (Permission, error) {
	resource, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		return Permission{}, err
	}
	return Permission{resource}, nil
}

func (s service) Create(ctx context.Context, org_id string, req CreateResourceRequest) (Permission, error) {
	if err := req.Validate(); err != nil {
		return Permission{}, err
	}
	id := entity.GenerateID()
	err := s.repo.Create(ctx, org_id, entity.Permission{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		return Permission{}, err
	}
	return s.Get(ctx, org_id, id)
}

func (s service) Update(ctx context.Context, org_id string, id string, req UpdateResourceRequest) (Permission, error) {
	if err := req.Validate(); err != nil {
		return Permission{}, err
	}

	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		return resource, err
	}
	resource.Name = req.Name
	if err := s.repo.Update(ctx, org_id, resource.Permission); err != nil {
		return resource, err
	}
	return resource, nil
}

func (s service) Delete(ctx context.Context, org_id string, id string) (Permission, error) {
	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		return Permission{}, err
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		return Permission{}, err
	}
	return resource, nil
}

func (s service) Query(ctx context.Context, org_id string) ([]Permission, error) {
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		return nil, err
	}
	result := []Permission{}
	for _, item := range items {
		result = append(result, Permission{item})
	}
	return result, nil
}