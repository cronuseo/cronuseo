package action

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, resource_id string, id string) (Action, error)
	Query(ctx context.Context, resource_id string) ([]Action, error)
	Create(ctx context.Context, resource_id string, input CreateActionRequest) (Action, error)
	Update(ctx context.Context, resource_id string, id string, input UpdateActionRequest) (Action, error)
	Delete(ctx context.Context, resource_id string, id string) (Action, error)
}

type Action struct {
	entity.Action
}

type CreateActionRequest struct {
	Key        string `json:"action_key" db:"action_key"`
	Name       string `json:"name" db:"name"`
	ResourceID string `json:"-" db:"resource_id"`
}

func (m CreateActionRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Key, validation.Required),
	)
}

type UpdateActionRequest struct {
	Name string `json:"name" db:"name"`
}

func (m UpdateActionRequest) Validate() error {
	return validation.ValidateStruct(&m)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) Get(ctx context.Context, resource_id string, id string) (Action, error) {
	resource, err := s.repo.Get(ctx, resource_id, id)
	if err != nil {
		return Action{}, &util.NotFoundError{Path: "Action"}
	}
	return Action{resource}, nil
}

func (s service) Create(ctx context.Context, resource_id string, req CreateActionRequest) (Action, error) {
	if err := req.Validate(); err != nil {
		return Action{}, &util.InvalidInputError{}
	}

	exists, _ := s.repo.ExistByKey(ctx, resource_id, req.Key)
	if exists {
		return Action{}, &util.AlreadyExistsError{Path: "Action"}
	}
	id := entity.GenerateID()
	err := s.repo.Create(ctx, resource_id, entity.Action{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		return Action{}, err
	}
	return s.Get(ctx, resource_id, id)
}

func (s service) Update(ctx context.Context, resource_id string, id string, req UpdateActionRequest) (Action, error) {
	if err := req.Validate(); err != nil {
		return Action{}, &util.InvalidInputError{}
	}

	resource, err := s.Get(ctx, resource_id, id)
	if err != nil {
		return Action{}, &util.NotFoundError{Path: "Action"}
	}
	resource.Name = req.Name
	if err := s.repo.Update(ctx, resource_id, resource.Action); err != nil {
		return resource, err
	}
	return resource, nil
}

func (s service) Delete(ctx context.Context, resource_id string, id string) (Action, error) {
	resource, err := s.Get(ctx, resource_id, id)
	if err != nil {
		return Action{}, &util.NotFoundError{Path: "Action"}
	}
	if err = s.repo.Delete(ctx, resource_id, id); err != nil {
		return Action{}, err
	}
	return resource, nil
}

func (s service) Query(ctx context.Context, resource_id string) ([]Action, error) {
	items, err := s.repo.Query(ctx, resource_id)
	if err != nil {
		return nil, err
	}
	result := []Action{}
	for _, item := range items {
		result = append(result, Action{item})
	}
	return result, nil
}
