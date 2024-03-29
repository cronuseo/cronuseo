package group

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (GroupResponse, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]Group, error)
	Create(ctx context.Context, org_id string, input CreateGroupRequest) (GroupResponse, error)
	Update(ctx context.Context, org_id string, id string, input UpdateGroupRequest) (GroupResponse, error)
	Delete(ctx context.Context, org_id string, id string) error
	Patch(ctx context.Context, org_id string, id string, input PatchGroupRequest) (GroupResponse, error)
}

type Group struct {
	mongo_entity.Group
}

type GroupResponse struct {
	ID          primitive.ObjectID            `json:"id" bson:"_id,omitempty"`
	Identifier  string                        `json:"identifier" bson:"identifier"`
	DisplayName string                        `json:"display_name" bson:"display_name"`
	Users       []mongo_entity.AssignedUser   `json:"users,omitempty" bson:"users"`
	Roles       []mongo_entity.AssignedRole   `json:"roles,omitempty" bson:"roles"`
	Policies    []mongo_entity.AssignedPolicy `json:"policies,omitempty" bson:"policies"`
}

type CreateGroupRequest struct {
	Identifier  string               `json:"identifier" bson:"identifier"`
	DisplayName string               `json:"display_name" bson:"display_name"`
	Roles       []primitive.ObjectID `json:"roles,omitempty" bson:"roles"`
	Users       []primitive.ObjectID `json:"users,omitempty" bson:"users"`
	Policies    []primitive.ObjectID `json:"policies,omitempty" bson:"policies"`
}

func (m CreateGroupRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Identifier, validation.Required),
	)
}

type UpdateGroupRequest struct {
	DisplayName *string `json:"display_name,omitempty" bson:"display_name"`
}

type PatchGroupRequest struct {
	AddedRoles      []primitive.ObjectID `json:"added_roles,omitempty" bson:"added_roles"`
	RemovedRoles    []primitive.ObjectID `json:"removed_roles,omitempty" bson:"removed_roles"`
	AddedUsers      []primitive.ObjectID `json:"added_users,omitempty" bson:"added_users"`
	RemovedUsers    []primitive.ObjectID `json:"removed_users,omitempty" bson:"removed_users"`
	AddedPolicies   []primitive.ObjectID `json:"added_policies,omitempty" bson:"added_policies"`
	RemovedPolicies []primitive.ObjectID `json:"removed_policies,omitempty" bson:"removed_policies"`
}

type UpdateGroup struct {
	DisplayName *string `json:"display_name,omitempty" bson:"display_name"`
}

type PatchGroup struct {
	AddedRoles      []primitive.ObjectID `json:"added_roles,omitempty" bson:"added_roles"`
	RemovedRoles    []primitive.ObjectID `json:"removed_roles,omitempty" bson:"removed_roles"`
	AddedUsers      []primitive.ObjectID `json:"added_users,omitempty" bson:"added_users"`
	RemovedUsers    []primitive.ObjectID `json:"removed_users,omitempty" bson:"removed_users"`
	AddedPolicies   []primitive.ObjectID `json:"added_policies,omitempty" bson:"added_policies"`
	RemovedPolicies []primitive.ObjectID `json:"removed_policies,omitempty" bson:"removed_policies"`
}

func (m UpdateGroupRequest) Validate() error {
	return validation.ValidateStruct(&m)
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}

// Get group by id.
func (s service) Get(ctx context.Context, org_id string, id string) (GroupResponse, error) {

	group, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the group.",
			zap.String("organization_id", org_id),
			zap.String("group_id", id))
		return GroupResponse{}, &util.NotFoundError{Path: "Group"}
	}
	return *group, err
}

// Create new group.
func (s service) Create(ctx context.Context, org_id string, req CreateGroupRequest) (GroupResponse, error) {

	// Validate group request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating group create request.")
		return GroupResponse{}, &util.InvalidInputError{Path: "Invalid input for group."}
	}

	// Check group already exists.
	exists, _ := s.repo.CheckGroupExistsByIdentifier(ctx, org_id, req.Identifier)
	if exists {
		s.logger.Debug("Group already exists.")
		return GroupResponse{}, &util.AlreadyExistsError{Path: "Group : " + req.Identifier}

	}

	// Generate group id.
	groupId := primitive.NewObjectID()

	for _, roleId := range req.Roles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}

	for _, userId := range req.Users {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid role id " + userId.String()}
		}
	}

	for _, policyId := range req.Policies {
		exists, _ := s.repo.CheckPolicyExistById(ctx, org_id, policyId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid policy id " + policyId.String()}
		}
	}

	var roles []primitive.ObjectID
	if req.Roles == nil {
		roles = []primitive.ObjectID{}
	} else {
		roles = req.Roles
	}

	var users []primitive.ObjectID
	if req.Users == nil {
		users = []primitive.ObjectID{}
	} else {
		users = req.Users
	}

	var policies []primitive.ObjectID
	if req.Policies == nil {
		policies = []primitive.ObjectID{}
	} else {
		policies = req.Policies
	}

	err := s.repo.Create(ctx, org_id, mongo_entity.Group{
		ID:          groupId,
		DisplayName: req.DisplayName,
		Identifier:  req.Identifier,
		Roles:       roles,
		Users:       users,
		Policies:    policies,
	})

	if err != nil {
		s.logger.Error("Error while creating group.",
			zap.String("organization_id", org_id))
		return GroupResponse{}, err
	}
	return s.Get(ctx, org_id, groupId.Hex())
}

// // Update group.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateGroupRequest) (GroupResponse, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Group not exists.", zap.String("group_id", id))
		return GroupResponse{}, &util.NotFoundError{Path: "Group " + id + " not exists."}
	}

	if err := s.repo.Update(ctx, org_id, id, UpdateGroup{
		DisplayName: req.DisplayName,
	}); err != nil {
		s.logger.Error("Error while updating group.",
			zap.String("organization_id", org_id),
			zap.String("group_id", id))
		return GroupResponse{}, err
	}
	return s.Get(ctx, org_id, id)
}

func (s service) Patch(ctx context.Context, org_id string, id string, req PatchGroupRequest) (GroupResponse, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("Group not exists.", zap.String("group_id", id))
		return GroupResponse{}, &util.NotFoundError{Path: "Group " + id + " not exists."}
	}

	// roles
	for _, roleId := range req.AddedRoles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}
	for _, roleId := range req.RemovedRoles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}
	added_roles := []primitive.ObjectID{}
	for _, roleId := range req.AddedRoles {
		already_added, _ := s.repo.CheckRoleAlreadyAssignToGroupById(ctx, org_id, id, roleId.Hex())
		if !already_added {
			added_roles = append(added_roles, roleId)
		}
	}

	removed_roles := []primitive.ObjectID{}
	for _, roleId := range req.RemovedRoles {
		already_added, _ := s.repo.CheckRoleAlreadyAssignToGroupById(ctx, org_id, id, roleId.Hex())
		if already_added {
			removed_roles = append(removed_roles, roleId)
		}
	}

	// users
	for _, userId := range req.AddedUsers {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid user id " + userId.String()}
		}
	}
	for _, userId := range req.RemovedUsers {
		exists, _ := s.repo.CheckUserExistById(ctx, org_id, userId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid user id " + userId.String()}
		}
	}
	added_users := []primitive.ObjectID{}
	for _, userId := range req.AddedUsers {
		already_added, _ := s.repo.CheckUserAlreadyAssignToGroupById(ctx, org_id, id, userId.Hex())
		if !already_added {
			added_users = append(added_users, userId)
		}
	}

	removed_users := []primitive.ObjectID{}
	for _, userId := range req.RemovedUsers {
		already_added, _ := s.repo.CheckUserAlreadyAssignToGroupById(ctx, org_id, id, userId.Hex())
		if already_added {
			removed_users = append(removed_users, userId)
		}
	}

	// policies
	for _, policyId := range req.AddedPolicies {
		exists, _ := s.repo.CheckPolicyExistById(ctx, org_id, policyId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid policy id " + policyId.String()}
		}
	}
	for _, policyId := range req.RemovedPolicies {
		exists, _ := s.repo.CheckPolicyExistById(ctx, org_id, policyId.Hex())
		if !exists {
			return GroupResponse{}, &util.InvalidInputError{Path: "Invalid policy id " + policyId.String()}
		}
	}
	added_policies := []primitive.ObjectID{}
	for _, policyId := range req.AddedPolicies {
		already_added, _ := s.repo.CheckPolicyAlreadyAssignToGroupById(ctx, org_id, id, policyId.Hex())
		if !already_added {
			added_policies = append(added_policies, policyId)
		}
	}

	removed_policies := []primitive.ObjectID{}
	for _, policyId := range req.RemovedPolicies {
		already_added, _ := s.repo.CheckPolicyAlreadyAssignToGroupById(ctx, org_id, id, policyId.Hex())
		if already_added {
			removed_policies = append(removed_policies, policyId)
		}
	}

	if err := s.repo.Patch(ctx, org_id, id, PatchGroup{
		AddedRoles:      added_roles,
		RemovedRoles:    removed_roles,
		AddedUsers:      added_users,
		RemovedUsers:    removed_users,
		AddedPolicies:   added_policies,
		RemovedPolicies: removed_policies,
	}); err != nil {
		s.logger.Error("Error while updating group.",
			zap.String("organization_id", org_id),
			zap.String("group_id", id))
		return GroupResponse{}, err
	}
	return s.Get(ctx, org_id, id)
}

// Delete group.
func (s service) Delete(ctx context.Context, org_id string, id string) error {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Group not exists.", zap.String("group_id", id))
		return &util.NotFoundError{Path: "Group " + id + " not exists."}

	}
	if err = s.repo.Delete(ctx, org_id, id); err != nil {
		s.logger.Error("Error while deleting group.",
			zap.String("organization_id", org_id),
			zap.String("group_id", id))
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

// // Get all group.
func (s service) Query(ctx context.Context, org_id string, filter Filter) ([]Group, error) {

	result := []Group{}
	items, err := s.repo.Query(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while retrieving all resources.",
			zap.String("organization_id", org_id))
		return []Group{}, err
	}

	for _, item := range *items {
		result = append(result, Group{item})
	}
	return result, err
}
