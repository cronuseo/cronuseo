package user

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, org_id string, id string) (User, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]User, error)
	Create(ctx context.Context, org_id string, input CreateUserRequest) (User, error)
	Update(ctx context.Context, org_id string, id string, input UpdateUserRequest) (User, error)
	Delete(ctx context.Context, org_id string, id string) error
	// Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) (User, error)
}

type User struct {
	mongo_entity.User
}

type CreateUserRequest struct {
	Username  string               `json:"username"`
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Roles     []primitive.ObjectID `json:"roles"`
}

func (m CreateUserRequest) Validate() error {

	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required),
	)
}

type UpdateUserRequest struct {
	FirstName    *string              `json:"first_name"`
	LastName     *string              `json:"last_name"`
	AddedRoles   []primitive.ObjectID `json:"added_roles"`
	RemovedRoles []primitive.ObjectID `json:"removed_roles"`
}

type UpdateUser struct {
	FirstName    *string              `json:"first_name"`
	LastName     *string              `json:"last_name"`
	AddedRoles   []primitive.ObjectID `json:"added_roles"`
	RemovedRoles []primitive.ObjectID `json:"removed_roles"`
}

func (m UpdateUserRequest) Validate() error {
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

// Get user by id.
func (s service) Get(ctx context.Context, org_id string, id string) (User, error) {

	user, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Error("Error while getting the user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return User{}, &util.NotFoundError{Path: "User"}
	}
	return User{*user}, err
}

// Create new user.
func (s service) Create(ctx context.Context, org_id string, req CreateUserRequest) (User, error) {

	// Validate user request.
	if err := req.Validate(); err != nil {
		s.logger.Error("Error while validating user create request.")
		return User{}, &util.InvalidInputError{Path: "Invalid input for user."}
	}

	// Check user already exists.
	exists, _ := s.repo.CheckUserExistsByIdentifier(ctx, org_id, req.Username)
	if exists {
		s.logger.Debug("User already exists.")
		return User{}, &util.AlreadyExistsError{Path: "User : " + req.Username}

	}
	// Generate user id.
	userId := primitive.NewObjectID()

	for _, roleId := range req.Roles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return User{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}
	// roles := []primitive.ObjectID{}
	// for _, roleId := range req.Roles {
	// 	already_added, _ := s.repo.CheckRoleAlreadyAssignToUserById(ctx, org_id, userId.Hex(), roleId.Hex())
	// 	if !already_added {
	// 		roles = append(roles, roleId)
	// 	}
	// }
	err := s.repo.Create(ctx, org_id, mongo_entity.User{
		ID:        userId,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Roles:     req.Roles,
	})

	if err != nil {
		s.logger.Error("Error while creating user.",
			zap.String("organization_id", org_id))
		return User{}, err
	}
	return s.Get(ctx, org_id, userId.Hex())
}

// // Update user.
func (s service) Update(ctx context.Context, org_id string, id string, req UpdateUserRequest) (User, error) {

	_, err := s.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return User{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}

	for _, roleId := range req.AddedRoles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return User{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
		}
	}
	for _, roleId := range req.RemovedRoles {
		exists, _ := s.repo.CheckRoleExistById(ctx, org_id, roleId.Hex())
		if !exists {
			return User{}, &util.InvalidInputError{Path: "Invalid role id " + roleId.String()}
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

	if err := s.repo.Update(ctx, org_id, id, UpdateUser{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		AddedRoles:   added_roles,
		RemovedRoles: removed_roles,
	}); err != nil {
		s.logger.Error("Error while updating user.",
			zap.String("organization_id", org_id),
			zap.String("user_id", id))
		return User{}, err
	}
	updatedUser, err := s.repo.Get(ctx, org_id, id)
	if err != nil {
		s.logger.Debug("User not exists.", zap.String("user_id", id))
		return User{}, &util.NotFoundError{Path: "User " + id + " not exists."}
	}
	return User{*updatedUser}, nil
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
		s.logger.Error("Error while retrieving all resources.",
			zap.String("organization_id", org_id))
		return []User{}, err
	}

	for _, item := range *items {
		result = append(result, User{item})
	}
	return result, err
}

// type UserPatchRequest struct {
// 	Operations []UserPatchOperation `json:"operations"`
// }

// type UserPatchOperation struct {
// 	Operation string  `json:"op"`
// 	Path      string  `json:"path"`
// 	Values    []Value `json:"values"`
// }

// type Value struct {
// 	Value string `json:"value"`
// }

// // Patch user. mainly patch user roles.
// func (s service) Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) (User, error) {

// 	user, err := s.Get(ctx, org_id, id)
// 	if err != nil {
// 		s.logger.Error("User not exists.", zap.String("user_id", id))
// 		return User{}, &util.NotFoundError{Path: "User " + id + " not exists."}

// 	}
// 	if err := s.repo.Patch(ctx, org_id, id, req); err != nil {
// 		s.logger.Error("Error while patching user.",
// 			zap.String("organization_id", org_id),
// 			zap.String("user_id", id),
// 		)
// 		return User{}, err
// 	}
// 	s.permissionCache.FlushAll(ctx)
// 	return user, nil
// }
