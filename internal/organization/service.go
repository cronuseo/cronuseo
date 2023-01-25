package organization

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, id string) (Organization, error)
	Query(ctx context.Context) ([]Organization, error)
	Create(ctx context.Context, input CreateOrganizationRequest) (Organization, error)
	Update(ctx context.Context, id string, input UpdateOrganizationRequest) (Organization, error)
	RefreshAPIKey(ctx context.Context, id string) (Organization, error)
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
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return service{repo: repo, logger: logger}
}

// Get organization by id.
func (s service) Get(ctx context.Context, id string) (Organization, error) {

	organization, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("Error while getting the organization.",
			zap.String("organization_id", id))
		return Organization{}, &util.NotFoundError{Path: "Organization"}
	}
	return Organization{organization}, nil
}

// Create new organization.
func (s service) Create(ctx context.Context, req CreateOrganizationRequest) (Organization, error) {

	// Validate organization
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating organization create request.")
		return Organization{}, &util.InvalidInputError{Path: "Invalid input for organization."}

	}

	// Check organization exists.
	exists, _ := s.repo.ExistByKey(ctx, req.Key)
	if exists {
		s.logger.Debug("Organization already exists.")
		return Organization{}, &util.AlreadyExistsError{Path: "Organization : " + req.Key + " already exists."}

	}
	// Generate organization id.
	id := entity.GenerateID()

	err := s.repo.Create(ctx, entity.Organization{
		ID:   id,
		Key:  req.Key,
		Name: req.Name,
	})
	if err != nil {
		s.logger.Error("Error while updating organization.",
			zap.String("organization_id", id))
		return Organization{}, err
	}
	return s.Get(ctx, id)
}

// Update organization by id.
func (s service) Update(ctx context.Context, id string, req UpdateOrganizationRequest) (Organization, error) {

	// Validate organization
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating organization create request.")
		return Organization{}, &util.InvalidInputError{Path: "Invalid input for organization."}

	}

	// Get organization
	organization, err := s.Get(ctx, id)
	if err != nil {
		s.logger.Debug("Organization not exists.", zap.String("organization_id", id))
		return Organization{}, &util.NotFoundError{Path: "Organization " + id + " not exists."}
	}

	organization.Name = req.Name
	if err := s.repo.Update(ctx, organization.Organization); err != nil {
		s.logger.Error("Error while updating organization.",
			zap.String("organization_id", id))
		return Organization{}, err
	}
	return organization, err
}

// Delete organization by id.
func (s service) Delete(ctx context.Context, id string) (Organization, error) {

	organization, err := s.Get(ctx, id)
	if err != nil {
		s.logger.Debug("Organization not exists.", zap.String("organization_id", id))
		return Organization{}, &util.NotFoundError{Path: "Organization " + id + " not exists."}
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Error while deleting organization.",
			zap.String("organization_id", id))
		return Organization{}, err
	}
	return organization, nil
}

// Refresh API key of the organization.
func (s service) RefreshAPIKey(ctx context.Context, id string) (Organization, error) {

	// Get organization
	organization, err := s.Get(ctx, id)
	if err != nil {
		s.logger.Debug("Organization not exists.", zap.String("organization_id", id))
		return Organization{}, &util.NotFoundError{Path: "Organization " + id + " not exists."}
	}

	// Generate new API key.
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return Organization{}, err
	}
	APIKey := base64.StdEncoding.EncodeToString(key)
	if err := s.repo.RefreshAPIKey(ctx, APIKey, organization.ID); err != nil {
		s.logger.Error("Error while updating organization.",
			zap.String("organization_id", id))
		return Organization{}, err
	}
	organization, err = s.Get(ctx, id)
	return organization, err
}

// Get all organizations.
func (s service) Query(ctx context.Context) ([]Organization, error) {

	items, err := s.repo.Query(ctx)
	if err != nil {
		s.logger.Error("Error while retrieving all organizations.")
		return []Organization{}, err
	}
	result := []Organization{}
	for _, item := range items {
		result = append(result, Organization{item})
	}
	return result, nil
}
