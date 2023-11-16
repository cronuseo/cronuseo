package role

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (Role, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Role, error)
	Create(ctx context.Context, org_id string, input CreateRoleRequest) (Role, error)
	Update(ctx context.Context, org_id string, id string, input UpdateRoleRequest) (Role, error)
	Patch(ctx context.Context, org_id string, id string, input PatchRoleRequest) (Role, error)
	Delete(ctx context.Context, org_id string, id string) error
	GetPermissions(ctx context.Context, org_id string, role_id string) ([]mongo_entity.Permission, error)
}

type Role struct {
	mongo_entity.Role
}

type CreateRoleRequest struct {
	Identifier  string                    `json:"identifier" bson:"identifier"`
	DisplayName string                    `json:"display_name" bson:"display_name"`
	Users       []primitive.ObjectID      `json:"users,omitempty" bson:"users"`
	Groups      []primitive.ObjectID      `json:"groups,omitempty" bson:"groups"`
	Permissions []mongo_entity.Permission `json:"permissions,omitempty" bson:"permissions"`
}

func (m CreateRoleRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Identifier, validation.Required),
	)
}

type UpdateRoleRequest struct {
	DisplayName *string `json:"display_name" bson:"display_name"`
}

type PatchRoleRequest struct {
	AddedUsers         []primitive.ObjectID      `json:"added_users,omitempty" bson:"added_users"`
	RemovedUsers       []primitive.ObjectID      `json:"removed_users,omitempty" bson:"removed_users"`
	AddedGroups        []primitive.ObjectID      `json:"added_groups,omitempty" bson:"added_groups"`
	RemovedGroups      []primitive.ObjectID      `json:"removed_groups,omitempty" bson:"removed_groups"`
	AddedPermissions   []mongo_entity.Permission `json:"added_permissions,omitempty" bson:"added_permissions"`
	RemovedPermissions []mongo_entity.Permission `json:"removed_permissions,omitempty" bson:"removed_permissions"`
}

type UpdateRole struct {
	DisplayName *string `json:"display_name" bson:"display_name"`
}

type PatchRole struct {
	AddedUsers         []primitive.ObjectID      `json:"added_users,omitempty" bson:"added_users"`
	RemovedUsers       []primitive.ObjectID      `json:"removed_users,omitempty" bson:"removed_users"`
	AddedGroups        []primitive.ObjectID      `json:"added_groups,omitempty" bson:"added_groups"`
	RemovedGroups      []primitive.ObjectID      `json:"removed_groups,omitempty" bson:"removed_groups"`
	AddedPermissions   []mongo_entity.Permission `json:"added_permissions,omitempty" bson:"added_permissions"`
	RemovedPermissions []mongo_entity.Permission `json:"removed_permissions,omitempty" bson:"removed_permissions"`
}

func (m UpdateRoleRequest) Validate() error {

	return validation.ValidateStruct(&m)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
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
	return Role{*role}, nil
}

// Create role.
func (s service) Create(ctx context.Context, org_id string, req CreateRoleRequest) (Role, error) {

	// Validate role request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating role creation request.")
		return Role{}, &util.InvalidInputError{Path: "Invalid input for role."}
	}

	exists, _ := s.repo.CheckRoleExistsByIdentifier(ctx, org_id, req.Identifier)
	if exists {
		s.logger.Debug("Role already exists.")
		return Role{}, &util.AlreadyExistsError{Path: "Role : " + req.Identifier + " already exists."}

	}

	// Generate role id.
	roleId := primitive.NewObjectID()

	for _, userId := range req.Users {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid user id " + userId.String()}
		}
	}

	for _, groupId := range req.Groups {
		exists, _ := s.repo.CheckGroupExistById(ctx, org_id, groupId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid group id " + groupId.String()}
		}
	}

	for _, permission := range req.Permissions {

		exists, _ := s.repo.CheckResourceActionExists(ctx, org_id, permission.Resource, permission.Action)
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid permission, Resource : " + permission.Resource + " Action : " + permission.Action}
		}
	}

	var users []primitive.ObjectID
	if req.Users == nil {
		users = []primitive.ObjectID{}
	} else {
		users = req.Users
	}

	var groups []primitive.ObjectID
	if req.Groups == nil {
		groups = []primitive.ObjectID{}
	} else {
		groups = req.Groups
	}

	var permisions []mongo_entity.Permission
	if req.Permissions == nil {
		permisions = []mongo_entity.Permission{}
	} else {
		permisions = req.Permissions
	}

	err := s.repo.Create(ctx, org_id, mongo_entity.Role{
		ID:          roleId,
		Identifier:  req.Identifier,
		DisplayName: req.DisplayName,
		Users:       users,
		Groups:      groups,
		Permissions: permisions,
	})

	if err != nil {
		s.logger.Error("Error while creating role.",
			zap.String("organization_id", org_id),
			zap.String("role identifier", req.Identifier))
		return Role{}, err
	}
	return s.Get(ctx, org_id, roleId.Hex())
}

// Update role.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateRoleRequest) (Role, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.", zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}

	if err := s.repo.Update(ctx, org_id, id, UpdateRole{
		DisplayName: req.DisplayName,
	}); err != nil {
		s.logger.Error("Error while updating role.", zap.String("organization_id", org_id), zap.String("role_id", id))
		return Role{}, err
	}
	updatedRole, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.", zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}
	return Role{*updatedRole}, nil
}

func (s service) Patch(ctx context.Context, org_id string, id string, req PatchRoleRequest) (Role, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.", zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}

	for _, userId := range req.AddedUsers {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid role id " + userId.String()}
		}
	}

	for _, userId := range req.RemovedUsers {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid role id " + userId.String()}
		}
	}

	for _, userId := range req.AddedUsers {
		already_added, _ := s.repo.CheckUserAlreadyAssignToRoleById(ctx, org_id, id, userId.Hex())
		if already_added {
			return Role{}, &util.InvalidInputError{Path: "Group : " + userId.Hex() + " already assigned to role :" + id}
		}
	}

	for _, userId := range req.RemovedUsers {
		already_added, _ := s.repo.CheckUserAlreadyAssignToRoleById(ctx, org_id, id, userId.Hex())
		if !already_added {
			return Role{}, &util.InvalidInputError{Path: "Group : " + userId.Hex() + " not assigned to role :" + id}
		}
	}

	// groups
	for _, groupId := range req.AddedGroups {
		exists, _ := s.repo.CheckGroupExistById(ctx, org_id, groupId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid group id " + groupId.String()}
		}
	}
	for _, groupId := range req.RemovedGroups {
		exists, _ := s.repo.CheckGroupExistById(ctx, org_id, groupId.Hex())
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid group id " + groupId.String()}
		}
	}

	for _, groupId := range req.AddedGroups {
		already_added, _ := s.repo.CheckGroupAlreadyAssignToRoleById(ctx, org_id, id, groupId.Hex())
		if already_added {
			return Role{}, &util.InvalidInputError{Path: "Group : " + groupId.Hex() + " already assigned to role :" + id}
		}
	}

	for _, groupId := range req.RemovedGroups {
		already_added, _ := s.repo.CheckGroupAlreadyAssignToRoleById(ctx, org_id, id, groupId.Hex())
		if !already_added {
			return Role{}, &util.InvalidInputError{Path: "Group : " + groupId.Hex() + " not assigned to role :" + id}
		}
	}

	// permissions
	for _, permission := range req.AddedPermissions {

		exists, _ := s.repo.CheckResourceActionExists(ctx, org_id, permission.Resource, permission.Action)
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid permission resource : " + permission.Resource + " action :" + permission.Action}
		}

		exists, _ = s.repo.CheckPermissionExists(ctx, org_id, id, permission.Resource, permission.Action)
		if exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid permission resource : " + permission.Resource + " action :" + permission.Action}
		}

	}

	for _, permission := range req.RemovedPermissions {

		exists, _ := s.repo.CheckResourceActionExists(ctx, org_id, permission.Resource, permission.Action)
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid permission resource : " + permission.Resource + " action :" + permission.Action}
		}

		exists, _ = s.repo.CheckPermissionExists(ctx, org_id, id, permission.Resource, permission.Action)
		if !exists {
			return Role{}, &util.InvalidInputError{Path: "Invalid permission resource : " + permission.Resource + " action :" + permission.Action}
		}

	}

	if err := s.repo.Patch(ctx, org_id, id, PatchRole{
		AddedUsers:         req.AddedUsers,
		RemovedUsers:        req.RemovedUsers,
		AddedGroups:        req.AddedGroups,
		RemovedGroups:      req.RemovedGroups,
		AddedPermissions:   req.AddedPermissions,
		RemovedPermissions: req.RemovedPermissions,
	}); err != nil {
		s.logger.Error("Error while updating role.", zap.String("organization_id", org_id), zap.String("role_id", id))
		return Role{}, err
	}
	updatedRole, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Role not exists.", zap.String("role_id", id))
		return Role{}, &util.NotFoundError{Path: "Role " + id + " not exists."}
	}
	return Role{*updatedRole}, nil
}

// Delete role.
func (s service) Delete(ctx context.Context, org_id string, id string) error {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Resource not exists.", zap.String("resource_id", id))
		return &util.NotFoundError{Path: "Resource " + id + " not exists."}
	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		s.logger.Error("Error while deleting resource.",
			zap.String("organization_id", org_id),
			zap.String("resource_id", id))
		return err
	}
	return nil
}

type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// Get all roles.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Role, error) {

	result := []Role{}
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while retrieving all resources.",
			zap.String("organization_id", org_id))
		return []Role{}, err
	}

	for _, item := range *items {
		result = append(result, Role{item})
	}
	return result, err
}

// Get permissions.
func (s service) GetPermissions(ctx context.Context, org_id string, role_id string) ([]mongo_entity.Permission, error) {

	result := []mongo_entity.Permission{}
	items, err := s.repo.GetPermissions(ctx, org_id, role_id)
	if err != nil {
		s.logger.Error("Error while retrieving all permission.",
			zap.String("organization_id", org_id), zap.String("role_id", role_id))
		return []mongo_entity.Permission{}, err
	}

	for _, item := range *items {
		result = append(result, mongo_entity.Permission{Action: item.Action, Resource: item.Resource})
	}
	return result, err
}
