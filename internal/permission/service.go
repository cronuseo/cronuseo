package permission

import (
	"context"
	"strconv"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

type Service interface {
	CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple, isCheck bool) (bool, error)
	DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckActions(ctx context.Context, org string, namespace string, request CheckActionsRequest) []string
	PatchPermissions(ctx context.Context, org string, namespace string, req PermissionPatchRequest) error
}

type Tuple struct {
	entity.Tuple
}

type service struct {
	repo            Repository
	permissionCache cache.PermissionCache
}

func NewService(repo Repository, cache cache.PermissionCache) Service {
	return service{repo: repo, permissionCache: cache}
}

func (s service) CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {

	exists, err := s.repo.CheckTuple(ctx, org, namespace, tuple)
	if exists {
		return &util.AlreadyExistsError{Path: "Tuple"}
	}
	if err != nil {
		return err
	}

	tuple = qualifiedTuple(org, tuple)
	return s.repo.CreateTuple(ctx, org, namespace, tuple)
}

func (s service) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple, isCheck bool) (bool, error) {

	tuple = qualifiedTuple(org, tuple)
	if isCheck {
		value, _ := s.permissionCache.Get(ctx, tuple)
		if value == "true" {
			return true, nil
		}
		if value == "false" {
			return false, nil
		}
	}
	allow, err := s.repo.CheckTuple(ctx, org, namespace, tuple)
	if err != nil {
		return false, err
	}
	b := strconv.FormatBool(allow)
	if isCheck {
		s.permissionCache.Set(ctx, tuple, b)
	}
	return allow, nil

}
func (s service) DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.DeleteTuple(ctx, org, namespace, tuple)
}

func (s service) CreatePermissions(ctx context.Context, org string, namespace string, permissions entity.CreatePermissionsRequest) error {
	for _, permission := range permissions.Permissions {
		tuple := entity.Tuple{Object: permission.Resource, Relation: permission.Action, SubjectId: permission.Role}
		err := s.CreateTuple(ctx, org, namespace, tuple)
		if err != nil {
			return err
		}
	}
	return nil
}

type PermissionPatchRequest struct {
	Operations []PermissionPatchOperation `json:"operations"`
}

type PermissionPatchOperation struct {
	Operation   string              `json:"op"`
	Permissions []entity.Permission `json:"permisssions"`
}

func (s service) PatchPermissions(ctx context.Context, org_id string, namespace string, req PermissionPatchRequest) error {

	org, err := s.repo.GetOrganization(ctx, org_id)
	if err != nil {
		return err
	}
	for _, operation := range req.Operations {
		switch operation.Operation {
		case "add":
			if len(operation.Permissions) > 0 {

				for _, permission := range operation.Permissions {
					tuple := entity.Tuple{Object: permission.Resource, Relation: permission.Action, SubjectId: permission.Role}
					exists, err := s.CheckTuple(ctx, org.Key, namespace, tuple, false)
					if exists {
						continue
					}
					if err != nil {
						return err
					}
					err = s.CreateTuple(ctx, org.Key, namespace, tuple)
					if err != nil {
						return err
					}
				}
			}
		case "remove":
			if len(operation.Permissions) > 0 {
				for _, permission := range operation.Permissions {
					tuple := entity.Tuple{Object: permission.Resource, Relation: permission.Action, SubjectId: permission.Role}
					exists, err := s.CheckTuple(ctx, org.Key, namespace, tuple, false)
					if !exists {
						continue
					}
					if err != nil {
						return err
					}
					err = s.DeleteTuple(ctx, org.Key, namespace, tuple)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	s.permissionCache.FlushAll(ctx)
	return nil
}

type CheckActionsRequest struct {
	Role     string   `json:"role"`
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

func (s service) CheckActions(ctx context.Context, org_id string, namespace string, request CheckActionsRequest) []string {

	org, err := s.repo.GetOrganization(ctx, org_id)
	if err != nil {
		return []string{}
	}
	allowed_actions := []string{}
	for _, action := range request.Actions {
		tuple := entity.Tuple{Object: request.Resource, Relation: action, SubjectId: request.Role}
		bool, _ := s.CheckTuple(ctx, org.Key, namespace, tuple, false)
		if bool {
			allowed_actions = append(allowed_actions, action)
		}
	}
	return allowed_actions
}

func qualifiedTuple(org string, tuple entity.Tuple) entity.Tuple {

	tuple.Object = org + "/" + tuple.Object
	tuple.SubjectId = org + "/" + tuple.SubjectId
	return tuple
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
