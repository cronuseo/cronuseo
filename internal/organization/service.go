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
	GetIdByIdentifier(ctx context.Context, identifier string) (string, error)
	Query(ctx context.Context) ([]Organization, error)
	Create(ctx context.Context, req OrganizationCreationRequest) (Organization, error)
	RegenerateAPIKey(ctx context.Context, id string) (Organization, error)
	Delete(ctx context.Context, id string) (Organization, error)
	CheckOrgExistByIdentifier(ctx context.Context, identifier string) (bool, error)
}

type Organization struct {
	mongo_entity.Organization
}

type OrganizationCreationRequest struct {
	Identifier  string                  `json:"identifier" bson:"identifier"`
	DisplayName string                  `json:"display_name" bson:"display_name"`
	API_KEY     string                  `json:"api_key" bson:"api_key"`
	Resources   []mongo_entity.Resource `json:"resources,omitempty" bson:"resources"`
	Users       []mongo_entity.User     `json:"users,omitempty" bson:"users"`
	Roles       []mongo_entity.Role     `json:"roles,omitempty" bson:"roles"`
	Groups      []mongo_entity.Group    `json:"groups,omitempty" bson:"groups"`
}

func (m OrganizationCreationRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Identifier, validation.Required),
		validation.Field(&m.DisplayName, validation.Required),
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

// Get organization id by identifier.
func (s service) GetIdByIdentifier(ctx context.Context, id string) (string, error) {

	orgId, err := s.repo.GetIdByIdentifier(ctx, id)
	if err != nil {
		return "", &util.NotFoundError{Path: "Organization"}
	}
	return orgId, nil
}

// Create new organization.
func (s service) Create(ctx context.Context, req OrganizationCreationRequest) (Organization, error) {

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

	// Generate API-Key for organization.
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return Organization{}, err

	}
	APIKey := base64.StdEncoding.EncodeToString(key)

	id, err := s.repo.Create(ctx, mongo_entity.Organization{
		Identifier:  req.Identifier,
		DisplayName: req.DisplayName,
		API_KEY:     APIKey,
		Users:       req.Users,
		Groups:      req.Groups,
		Roles:       req.Roles,
		Resources:   req.Resources,
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

// Regenerate API key of the organization.
func (s service) RegenerateAPIKey(ctx context.Context, id string) (Organization, error) {

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

func (s service) CheckOrgExistByIdentifier(ctx context.Context, identifier string) (bool, error) {

	return s.repo.CheckOrgExistByIdentifier(ctx, identifier)
}
