package organization

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, id string) (Organization, error)
	Query(ctx context.Context) ([]Organization, error)
	Create(ctx context.Context, input CreateOrganizationRequest) (Organization, error)
	RefreshAPIKey(ctx context.Context, id string) (Organization, error)
	Delete(ctx context.Context, id string) (Organization, error)
}

type Organization struct {
	mongo_entity.Organization
}

type CreateOrganizationRequest struct {
	Identifier  string `json:"identifier" db:"identifier"`
	DisplayName string `json:"display_name" db:"display_name"`
}

func (m CreateOrganizationRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Identifier, validation.Required),
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

	org, err := s.repo.Get(ctx, id)
	if err != nil {
		return Organization{}, &util.NotFoundError{Path: "Organization"}
	}
	return Organization{*org}, nil
}

// Create new organization.
func (s service) Create(ctx context.Context, req CreateOrganizationRequest) (Organization, error) {

	// Validate organization
	if err := req.Validate(); err != nil {
		return Organization{}, &util.InvalidInputError{Path: "Invalid input for organization."}

	}

	// Check organization exists.
	exists, _ := s.repo.CheckOrgExistByIdentifier(ctx, req.Identifier)
	if exists {
		s.logger.Debug("Organization already exists.")
		return Organization{}, &util.AlreadyExistsError{Path: "Organization : " + req.Identifier + " already exists."}

	}

	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return Organization{}, err

	}
	APIKey := base64.StdEncoding.EncodeToString(key)

	id, err := s.repo.Create(ctx, mongo_entity.Organization{
		Identifier:  req.Identifier,
		DisplayName: req.DisplayName,
		API_KEY:     APIKey,
	})
	if err != nil {
		s.logger.Error("Error while creating organization.")
		return Organization{}, err
	}
	return s.Get(ctx, id)
}

// Delete organization by id.
func (s service) Delete(ctx context.Context, id string) (Organization, error) {

	organization, err := s.Get(ctx, id)
	if err != nil {
		s.logger.Debug("Organization not exists.", zap.String("organization_id", id))
		return Organization{}, &util.NotFoundError{Path: "Organization " + id + " not exists."}
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Error while deleting organization.", zap.String("organization_id", id))
		return Organization{}, err
	}
	return organization, nil
}

// Refresh API key of the organization.
func (s service) RefreshAPIKey(ctx context.Context, id string) (Organization, error) {

	// Get organization
	exists, _ := s.repo.CheckOrgExistById(ctx, id)
	if !exists {
		s.logger.Debug("Organization not exists.", zap.String("organization_id", id))
		return Organization{}, &util.NotFoundError{Path: "Organization " + id + " not exists."}
	}

	// Generate new API key.
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return Organization{}, err
	}
	APIKey := base64.StdEncoding.EncodeToString(key)
	if err := s.repo.RefreshAPIKey(ctx, APIKey, id); err != nil {
		s.logger.Error("Error while updating organization.", zap.String("organization_id", id))
		return Organization{}, err
	}
	organization, err := s.Get(ctx, id)
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
