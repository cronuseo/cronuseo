package role

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Role, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Role, error)
	QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]Role, error)
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
	repo            Repository
	permissionCache cache.PermissionCache
	logger          *zap.Logger
}

func NewService(repo Repository, permissionCache cache.PermissionCache, logger *zap.Logger) Service {

	return service{repo: repo, permissionCache: permissionCache, logger: logger}
}

// Get role by id.
func (s service) Get(ctx context.Context, org_id string, id string) (Role, error) {

	role, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the role.",
			zap.String("organization_id", org_id),
			zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role"}
	}
	return Role{role}, nil
}

// Create role.
func (s service) Create(ctx context.Context, org_id string, req CreateRoleRequest) (Role, error) {

	// Validate role request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating role request.")
		return Role{}, &util.InvalidInputError{Path: "Invalid input for role."}
	}

	exists, _ := s.repo.ExistByKey(ctx, req.Key)
	if exists {
		s.logger.Debug("Role already exists.")
		return Role{}, &util.AlreadyExistsError{Path: "Role : " + req.Key + " already exists."}

	}

	// Generate role id.
	id := entity.GenerateID()

	err := s.repo.Create(ctx, org_id, entity.Role{
		ID:    id,
		Key:   req.Key,
		Name:  req.Name,
		Users: req.Users,
	})

	if err != nil {
		s.logger.Error("Error while creating role.",
			zap.String("organization_id", org_id),
			zap.String("role key", req.Key))
		return Role{}, err
	}
	return s.Get(ctx, org_id, id)
}

// Update role.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateRoleRequest) (Role, error) {

	// Validate role request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating role request.")
		return Role{}, &util.InvalidInputError{Path: "Invalid input for role."}
	}

	role, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.")
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}
	role.Name = req.Name
	if err := s.repo.Update(ctx, org_id, role.Role); err != nil {
		s.logger.Error("Error while creating role.",
			zap.String("organization_id", org_id),
			zap.String("role_id", id))
		return Role{}, err
	}
	return role, err
}

// Delete role.
func (s service) Delete(ctx context.Context, org_id string, id string) (Role, error) {

	role, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Role not exists.")
		return role, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}
	if err = s.repo.Delete(ctx, role.Role); err != nil {
		s.logger.Error("Error while deleting role.",
			zap.String("organization_id", org_id),
			zap.String("role_id", id))
		return Role{}, err
	}
	s.permissionCache.FlushAll(ctx)
	return role, err
}

type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// Get all roles.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Role, error) {

	items, err := s.repo.Query(ctx, org_id, filter)
	if err != nil {
		s.logger.Error("Error while retrieving all roles.",
			zap.String("organization_id", org_id))
		return []Role{}, err
	}
	result := []Role{}
	for _, item := range items {
		result = append(result, Role{item})
	}
	return result, err
}

// Get all roles by user id.
func (s service) QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]Role, error) {

	items, err := s.repo.QueryByUserID(ctx, org_id, user_id, filter)
	if err != nil {
		s.logger.Error("Error while retrieving all roles.",
			zap.String("organization_id", org_id),
			zap.String("user_id", user_id))
		return nil, err
	}
	result := []Role{}
	for _, item := range items {
		result = append(result, Role{item})
	}
	return result, err
}
