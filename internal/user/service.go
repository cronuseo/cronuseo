package user

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (UserResponse, error)
	GetIdByIdentifier(ctx context.Context, org_id string, identifier string) (string, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]User, error)
	Create(ctx context.Context, org_id string, input CreateUserRequest) (UserResponse, error)
	Sync(ctx context.Context, org_id string, input SyncUserRequest) (UserResponse, error)
	Update(ctx context.Context, org_id string, id string, input UpdateUserRequest) (UserResponse, error)
	Patch(ctx context.Context, org_id string, id string, input PatchUserRequest) (UserResponse, error)
	Delete(ctx context.Context, org_id string, id string) error
}

type User struct {
	mongo_entity.User
}

type UserResponse struct {
	ID             primitive.ObjectID            `json:"id" bson:"_id,omitempty"`
	Username       string                        `json:"username" bson:"username"`
	Identifier     string                        `json:"identifier" bson:"identifier"`
	UserProperties map[string]interface{}        `json:"user_properties" bson:"user_properties"`
	Roles          []mongo_entity.AssignedRole   `json:"roles,omitempty" bson:"roles"`
	Groups         []mongo_entity.AssignedGroup  `json:"groups,omitempty" bson:"groups"`
	Policies       []mongo_entity.AssignedPolicy `json:"policies,omitempty" bson:"policies"`
}

type CreateUserRequest struct {
	Username       string                 `json:"username" bson:"username"`
	Identifier     string                 `json:"identifier" bson:"identifier"`
	UserProperties map[string]interface{} `json:"user_properties" bson:"user_properties"`
	Roles          []primitive.ObjectID   `json:"roles,omitempty" bson:"roles"`
	Groups         []primitive.ObjectID   `json:"groups,omitempty" bson:"groups"`
	Policies       []primitive.ObjectID   `json:"policies,omitempty" bson:"policies"`
}

type SyncUserRequest struct {
	Username       string                 `json:"username" bson:"username"`
	Identifier     string                 `json:"identifier" bson:"identifier"`
	UserProperties map[string]interface{} `json:"user_properties" bson:"user_properties"`
	Roles          []string               `json:"roles,omitempty" bson:"roles"`
	Groups         []string               `json:"groups,omitempty" bson:"groups"`
}

func (m CreateUserRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Identifier, validation.Required),
	)
}

func (m SyncUserRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required),
		validation.Field(&m.Identifier, validation.Required),
	)
}

type UpdateUserRequest struct {
	UserProperties map[string]interface{} `json:"user_properties" bson:"user_properties"`
}

type PatchUserRequest struct {
	UserProperties  map[string]interface{} `json:"user_properties,omitempty" bson:"user_properties"`
	AddedRoles      []primitive.ObjectID   `json:"added_roles,omitempty" bson:"added_roles"`
	RemovedRoles    []primitive.ObjectID   `json:"removed_roles,omitempty" bson:"removed_roles"`
	AddedGroups     []primitive.ObjectID   `json:"added_groups,omitempty" bson:"added_groups"`
	RemovedGroups   []primitive.ObjectID   `json:"removed_groups,omitempty" bson:"removed_groups"`
	AddedPolicies   []primitive.ObjectID   `json:"added_policies,omitempty" bson:"added_policies"`
	RemovedPolicies []primitive.ObjectID   `json:"removed_policies,omitempty" bson:"removed_policies"`
}

type UpdateUser struct {
	UserProperties map[string]interface{} `json:"user_properties" bson:"user_properties"`
}

type PatchUser struct {
	UserProperties  map[string]interface{} `json:"user_properties,omitempty" bson:"user_properties"`
	AddedRoles      []primitive.ObjectID   `json:"added_roles,omitempty" bson:"added_roles"`
	RemovedRoles    []primitive.ObjectID   `json:"removed_roles,omitempty" bson:"removed_roles"`
	AddedGroups     []primitive.ObjectID   `json:"added_groups,omitempty" bson:"added_groups"`
	RemovedGroups   []primitive.ObjectID   `json:"removed_groups,omitempty" bson:"removed_groups"`
	AddedPolicies   []primitive.ObjectID   `json:"added_policies,omitempty" bson:"added_policies"`
	RemovedPolicies []primitive.ObjectID   `json:"removed_policies,omitempty" bson:"removed_policies"`
}

func (m UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&m)
}

type service struct {
	repo        Repository
	logger      *zap.Logger
	roleService role.Service
}

func NewService(repo Repository, logger *zap.Logger, roleService role.Service) Service {

	return service{repo: repo, logger: logger, roleService: roleService}
}

// Get user by id.
func (s service) Get(ctx context.Context, org_id string, id string) (UserResponse, error) {

	user, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return UserResponse{}, &util.NotFoundError{Path: "User"}
	}
	return *user, err
}

// Get user id by identifier.
func (s service) GetIdByIdentifier(ctx context.Context, org_id string, identifier string) (string, error) {

	userId, err := s.repo.GetIdByIdentifier(ctx, org_id, identifier)
	if err != nil {
		s.logger.Error("Error while getting the user id.",
			zap.String("organization_id", org_id),
			zap.String("identifier", identifier))
		return "", &util.NotFoundError{Path: "User"}
	}
	return userId, err
}

// Create new user.
func (s service) Create(ctx context.Context, org_id string, req CreateUserRequest) (UserResponse, error) {

	// Validate user request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating user create request.")
		return UserResponse{}, &util.InvalidInputError{Path: "Invalid input for user."}
	}

	// Check user already exists.
	exists, _ := s.repo.CheckUserExistsByIdentifier(ctx, org_id, req.Username)
	if exists {
		s.logger.Debug("User already exists.")
		return UserResponse{}, &util.AlreadyExistsError{Path: "User : " + req.Username}

	}

	// Generate user id.
	userId := primitive.NewObjectID()

	for _, roleId := range req.Roles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}

	for _, groupId := range req.Groups {
		exists, _ := s.repo.CheckGroupExistById(ctx, org_id, groupId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid group id " + groupId.String()}
		}
	}

	for _, policyId := range req.Policies {
		exists, _ := s.repo.CheckPolicyExistById(ctx, org_id, policyId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid policy id " + policyId.String()}
		}
	}

	var roles []primitive.ObjectID
	if req.Roles == nil {
		roles = []primitive.ObjectID{}
	} else {
		roles = req.Roles
	}

	var groups []primitive.ObjectID
	if req.Groups == nil {
		groups = []primitive.ObjectID{}
	} else {
		groups = req.Groups
	}

	var policies []primitive.ObjectID
	if req.Policies == nil {
		policies = []primitive.ObjectID{}
	} else {
		policies = req.Policies
	}

	err := s.repo.Create(ctx, org_id, mongo_entity.User{
		ID:             userId,
		Username:       req.Username,
		Identifier:     req.Identifier,
		UserProperties: req.UserProperties,
		Roles:          roles,
		Groups:         groups,
		Policies:       policies,
	})

	if err != nil {
		s.logger.Error("Error while creating user.",
			zap.String("organization_id", org_id))
		return UserResponse{}, err
	}
	return s.Get(ctx, org_id, userId.Hex())
}

// Sync user.
func (s service) Sync(ctx context.Context, org_identifier string, req SyncUserRequest) (UserResponse, error) {

	// Validate user request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating user create request.")
		return UserResponse{}, &util.InvalidInputError{Path: "Invalid input for user."}
	}

	org_id, err := s.repo.GetOrgIdByIdentifier(ctx, org_identifier)
	if err != nil {
		s.logger.Error("Error while syncing user. Invalid org identifier", zap.String("organization_identifier", org_identifier))
		return UserResponse{}, err
	}

	// Check user already exists.
	exists, _ := s.repo.CheckUserExistsByIdentifier(ctx, org_id, req.Identifier)
	roleIds := []primitive.ObjectID{}
	if req.Roles != nil && len(req.Roles) > 0 {
		for _, roleIdentifier := range req.Roles {
			exists, _ := s.roleService.CheckRoleExistsByIdentifier(ctx, org_id, roleIdentifier)
			if !exists {
				roleCreateRequest := role.CreateRoleRequest{
					Identifier:  roleIdentifier,
					DisplayName: roleIdentifier,
				}
				role, _ := s.roleService.Create(ctx, org_id, roleCreateRequest)
				roleIds = append(roleIds, role.ID)
			} else {
				role, _ := s.roleService.GetRoleByIdentifier(ctx, org_id, roleIdentifier)
				roleIds = append(roleIds, role.ID)
			}
		}
	}
	if exists {
		s.logger.Debug("User already exists.")
		id, _ := s.GetIdByIdentifier(ctx, org_id, req.Identifier)
		addedRoles := []primitive.ObjectID{}
		for _, roleId := range roleIds {
			already_added, _ := s.repo.CheckRoleAlreadyAssignToUserById(ctx, org_id, id, roleId.Hex())
			if !already_added {
				addedRoles = append(addedRoles, roleId)
			}
		}
		if len(addedRoles) > 0 {
			patchUserRequest := PatchUserRequest{
				AddedRoles: addedRoles,
			}
			s.Patch(ctx, org_id, id, patchUserRequest)
		}
		return s.Get(ctx, org_id, id)

	} else {
		// Generate user id.
		userId := primitive.NewObjectID()

		err = s.repo.Create(ctx, org_id, mongo_entity.User{
			ID:         userId,
			Username:   req.Username,
			Identifier: req.Identifier,
			Roles:      roleIds,
			Groups:     []primitive.ObjectID{},
			Policies:   []primitive.ObjectID{},
		})

		if err != nil {
			s.logger.Error("Error while syncing user.",
				zap.String("organization_id", org_id))
			return UserResponse{}, err
		}
		return s.Get(ctx, org_id, userId.Hex())
	}
}

// // Update user.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateUserRequest) (UserResponse, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return UserResponse{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}

	if err := s.repo.Update(ctx, org_id, id, UpdateUser{
		UserProperties: req.UserProperties,
	}); err != nil {
		s.logger.Error("Error while updating user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return UserResponse{}, err
	}
	return s.Get(ctx, org_id, id)
}

func (s service) Patch(ctx context.Context, org_id string, id string, req PatchUserRequest) (UserResponse, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return UserResponse{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}

	// roles
	for _, roleId := range req.AddedRoles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}
	for _, roleId := range req.RemovedRoles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}
	added_roles := []primitive.ObjectID{}
	for _, roleId := range req.AddedRoles {
		already_added, _ := s.repo.CheckRoleAlreadyAssignToUserById(ctx, org_id, id, roleId.Hex())
		if !already_added {
			added_roles = append(added_roles, roleId)
		}
	}

	removed_roles := []primitive.ObjectID{}
	for _, roleId := range req.RemovedRoles {
		already_added, _ := s.repo.CheckRoleAlreadyAssignToUserById(ctx, org_id, id, roleId.Hex())
		if already_added {
			removed_roles = append(removed_roles, roleId)
		}
	}

	// groups
	for _, groupId := range req.AddedGroups {
		exists, _ := s.repo.CheckGroupExistById(ctx, org_id, groupId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid group id " + groupId.String()}
		}
	}
	for _, groupId := range req.RemovedGroups {
		exists, _ := s.repo.CheckGroupExistById(ctx, org_id, groupId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid group id " + groupId.String()}
		}
	}
	added_groups := []primitive.ObjectID{}
	for _, groupId := range req.AddedGroups {
		already_added, _ := s.repo.CheckGroupAlreadyAssignToUserById(ctx, org_id, id, groupId.Hex())
		if !already_added {
			added_groups = append(added_groups, groupId)
		}
	}

	removed_groups := []primitive.ObjectID{}
	for _, groupId := range req.RemovedGroups {
		already_added, _ := s.repo.CheckGroupAlreadyAssignToUserById(ctx, org_id, id, groupId.Hex())
		if already_added {
			removed_groups = append(removed_groups, groupId)
		}
	}

	// policies
	for _, policyId := range req.AddedPolicies {
		exists, _ := s.repo.CheckPolicyExistById(ctx, org_id, policyId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid policy id " + policyId.String()}
		}
	}
	for _, policyId := range req.RemovedPolicies {
		exists, _ := s.repo.CheckPolicyExistById(ctx, org_id, policyId.Hex())
		if !exists {
			return UserResponse{}, &util.InvalidInputError{Path: "Invalid policy id " + policyId.String()}
		}
	}
	added_policies := []primitive.ObjectID{}
	for _, policyId := range req.AddedPolicies {
		already_added, _ := s.repo.CheckPolicyAlreadyAssignToUserById(ctx, org_id, id, policyId.Hex())
		if !already_added {
			added_policies = append(added_policies, policyId)
		}
	}

	removed_policies := []primitive.ObjectID{}
	for _, policyId := range req.RemovedPolicies {
		already_added, _ := s.repo.CheckPolicyAlreadyAssignToUserById(ctx, org_id, id, policyId.Hex())
		if already_added {
			removed_policies = append(removed_policies, policyId)
		}
	}

	if err := s.repo.Patch(ctx, org_id, id, PatchUser{
		UserProperties:  req.UserProperties,
		AddedRoles:      added_roles,
		RemovedRoles:    removed_roles,
		AddedGroups:     added_groups,
		RemovedGroups:   removed_groups,
		AddedPolicies:   added_policies,
		RemovedPolicies: removed_policies,
	}); err != nil {
		s.logger.Error("Error while updating user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return UserResponse{}, err
	}
	return s.Get(ctx, org_id, id)
}

// Delete user.
func (s service) Delete(ctx context.Context, org_id string, id string) error {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("User not exists.", zap.String("user_id", id))
		return &util.NotFoundError{Path: "User " + id + " not exists."}

	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		s.logger.Error("Error while deleting user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return err
	}
	return nil
}

// Pagination filter.
type Filter struct {
	Cursor int    `json:"cursor" query:"cursor"`
	Limit  int    `json:"limit" query:"limit"`
	Name   string `json:"name" query:"name"`
}

// // Get all user.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]User, error) {

	result := []User{}
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while retrieving all user.",
			zap.String("organization_id", org_id))
		return []User{}, err
	}

	for _, item := range *items {
		result = append(result, User{item})
	}
	return result, err
}
