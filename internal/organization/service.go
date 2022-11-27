package organization

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, id string) (Organization, error)
	Query(ctx context.Context) ([]Organization, error)
	Create(ctx context.Context, input CreateOrganizationRequest) (Organization, error)
	Update(ctx context.Context, id string, input UpdateOrganizationRequest) (Organization, error)
	Delete(ctx context.Context, id string) (Organization, error)
}

type Organization struct {
	entity.Organization
}

type CreateOrganizationRequest struct {
	Key  string `json:"org_key" db:"org_key"`
	Name string `json:"name" db:"name"`
}

func (m CreateOrganizationRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Key, validation.Required),
	)
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" db:"name"`
}

func (m UpdateOrganizationRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) Get(ctx context.Context, id string) (Organization, error) {
	organization, err := s.repo.Get(ctx, id)
	if err != nil {
		return Organization{}, &util.NotFoundError{Path: "Organization"}
	}
	return Organization{organization}, nil
}

func (s service) Create(ctx context.Context, req CreateOrganizationRequest) (Organization, error) {

	//validate organization
	if err := req.Validate(); err != nil {
		return Organization{}, &util.InvalidInputError{}
	}

	//check organixation exists
	exists, _ := s.repo.ExistByKey(ctx, req.Key)
	if exists {
		return Organization{}, &util.AlreadyExistsError{Path: "Organization"}
	}

	id := entity.GenerateID()
	err := s.repo.Create(ctx, entity.Organization{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		return Organization{}, err
	}
	return s.Get(ctx, id)
}

func (s service) Update(ctx context.Context, id string, req UpdateOrganizationRequest) (Organization, error) {

	//validate organization
	if err := req.Validate(); err != nil {
		return Organization{}, &util.InvalidInputError{}
	}

	organization, err := s.Get(ctx, id)
	if err != nil {
		return Organization{}, &util.NotFoundError{Path: "Organization"}
	}
	organization.Name = req.Name

	if err := s.repo.Update(ctx, organization.Organization); err != nil {
		return organization, err
	}
	return organization, nil
}

func (s service) Delete(ctx context.Context, id string) (Organization, error) {
	organization, err := s.Get(ctx, id)
	if err != nil {
		return Organization{}, &util.NotFoundError{Path: "Organization"}
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Organization{}, err
	}
	return organization, nil
}

func (s service) Query(ctx context.Context) ([]Organization, error) {
	items, err := s.repo.Query(ctx)
	if err != nil {
		return nil, err
	}
	result := []Organization{}
	for _, item := range items {
		result = append(result, Organization{item})
	}
	return result, nil
}
