package role

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Role, error)
	Query(ctx context.Context, org_id string) ([]Role, error)
	Create(ctx context.Context, org_id string, input CreateRoleRequest) (Role, error)
	Update(ctx context.Context, org_id string, id string, input UpdateRoleRequest) (Role, error)
	Delete(ctx context.Context, org_id string, id string) (Role, error)
}

type Role struct {
	entity.Role
}

type CreateRoleRequest struct {
	Key   string          `json:"role_key" db:"role_key"`
	Name  string          `json:"name" db:"name"`
	OrgID string          `json:"-" db:"org_id"`
	Users []entity.UserID `json:"users"`
}

func (m CreateRoleRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Key, validation.Required),
	)
}

type UpdateRoleRequest struct {
	Name string `json:"name" db:"name"`
}

func (m UpdateRoleRequest) Validate() error {
	return validation.ValidateStruct(&m)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) Get(ctx context.Context, org_id string, id string) (Role, error) {
	role, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		return Role{}, err
	}
	return Role{role}, nil
}

func (s service) Create(ctx context.Context, org_id string, req CreateRoleRequest) (Role, error) {
	if err := req.Validate(); err != nil {
		return Role{}, err
	}
	id := entity.GenerateID()
	err := s.repo.Create(ctx, org_id, entity.Role{
		ID:    id,
		Key:   req.Key,
		Name:  req.Name,
		Users: req.Users,
	})
	if err != nil {
		return Role{}, err
	}
	return s.Get(ctx, org_id, id)
}

func (s service) Update(ctx context.Context, org_id string, id string, req UpdateRoleRequest) (Role, error) {
	if err := req.Validate(); err != nil {
		return Role{}, err
	}

	role, err := s.Get(ctx, org_id, id)
	if err != nil {
		return role, err
	}
	role.Name = req.Name
	if err := s.repo.Update(ctx, org_id, role.Role); err != nil {
		return role, err
	}
	return role, nil
}

func (s service) Delete(ctx context.Context, org_id string, id string) (Role, error) {
	role, err := s.Get(ctx, org_id, id)
	if err != nil {
		return Role{}, err
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		return Role{}, err
	}
	return role, nil
}

func (s service) Query(ctx context.Context, org_id string) ([]Role, error) {
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		return nil, err
	}
	result := []Role{}
	for _, item := range items {
		result = append(result, Role{item})
	}
	return result, nil
}
