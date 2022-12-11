package resource

import (
	"context"
	"log"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Resource, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Resource, error)
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
		return Resource{}, &util.NotFoundError{Path: "Organization"}
	}
	return Resource{resource}, nil
}

func (s service) Create(ctx context.Context, org_id string, req CreateResourceRequest) (Resource, error) {

	//validate organization
	if err := req.Validate(); err != nil {
		return Resource{}, &util.InvalidInputError{}
	}

	//check organization exists
	exists, _ := s.repo.ExistByKey(ctx, req.Key)
	if exists {
		return Resource{}, &util.AlreadyExistsError{Path: "Resource"}
	}

	id := entity.GenerateID()
	err := s.repo.Create(ctx, org_id, entity.Resource{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		log.Println(err.Error())
		return Resource{}, err
	}
	return s.Get(ctx, org_id, id)
}

func (s service) Update(ctx context.Context, org_id string, id string, req UpdateResourceRequest) (Resource, error) {
	if err := req.Validate(); err != nil {
		return Resource{}, &util.InvalidInputError{}
	}

	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		return resource, &util.NotFoundError{Path: "Resource"}
	}
	resource.Name = req.Name
	if err := s.repo.Update(ctx, org_id, resource.Resource); err != nil {
		log.Println(err.Error())
		return resource, err
	}
	return resource, nil
}

func (s service) Delete(ctx context.Context, org_id string, id string) (Resource, error) {
	resource, err := s.Get(ctx, org_id, id)
	if err != nil {
		return Resource{}, &util.NotFoundError{Path: "Resource"}
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		log.Println(err.Error())
		return Resource{}, err
	}
	return resource, nil
}

type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Resource, error) {

	if filter.Limit == 0 {
		filter.Limit = 10
	}
	items, err := s.repo.Query(ctx, org_id, filter)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	result := []Resource{}
	for _, item := range items {
		result = append(result, Resource{item})
	}
	return result, nil
}
