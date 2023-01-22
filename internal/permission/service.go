package permission

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"
)

// Permission service handle all permission related operations. basically its handle all CRUD operations on keto.
type Service interface {
	CreateTuple(ctx context.Context, org string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, tuple entity.Tuple) (bool, error)
	DeleteTuple(ctx context.Context, org string, tuple entity.Tuple) error
	CheckActions(ctx context.Context, org string, request CheckActionsRequest) ([]string, error)
	PatchPermissions(ctx context.Context, org string, req PermissionPatchRequest) error
}

type Tuple struct {
	entity.Tuple
}

type service struct {
	repo            Repository
	permissionCache cache.PermissionCache
	logger          *zap.Logger
}

func NewService(repo Repository, cache cache.PermissionCache, logger *zap.Logger) Service {

	return service{
		repo:            repo,
		permissionCache: cache,
		logger:          logger,
	}
}

// Create tuple in keto.
func (s service) CreateTuple(ctx context.Context, org string, tuple entity.Tuple) error {

	qTuple := qualifiedTuple(org, tuple)
	err := s.repo.CreateTuple(ctx, org, qTuple)
	if err != nil {
		s.logger.Error("Error while creating tuple with keto.",
			zap.String("subject", tuple.SubjectId),
			zap.String("object", tuple.Object),
			zap.String("relation", tuple.Relation),
		)
	}
	return err
}

// Check tuple in keto.
func (s service) CheckTuple(ctx context.Context, org string, tuple entity.Tuple) (bool, error) {

	qTuple := qualifiedTuple(org, tuple)
	allow, err := s.repo.CheckTuple(ctx, org, qTuple)
	if err != nil {
		s.logger.Error("Error while checking tuple with keto.",
			zap.String("subject", tuple.SubjectId),
			zap.String("object", tuple.Object),
			zap.String("relation", tuple.Relation),
		)
		return false, err
	}
	return allow, nil

}

// Delete tuple in keto.
func (s service) DeleteTuple(ctx context.Context, org string, tuple entity.Tuple) error {

	qTuple := qualifiedTuple(org, tuple)
	err := s.repo.DeleteTuple(ctx, org, qTuple)
	if err != nil {
		s.logger.Error("Error while deleting tuple with keto.",
			zap.String("subject", tuple.SubjectId),
			zap.String("object", tuple.Object),
			zap.String("relation", tuple.Relation),
		)
	}
	return err
}

type PermissionPatchRequest struct {
	Operations []PermissionPatchOperation `json:"operations"`
}

type PermissionPatchOperation struct {
	Operation   string              `json:"op"`
	Permissions []entity.Permission `json:"permissions"`
}

// Add or remove permissions in keto.
func (s service) PatchPermissions(ctx context.Context, org_id string, req PermissionPatchRequest) error {

	org, err := s.repo.GetOrganization(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while getting organization from database.",
			zap.String("organization_id", org_id),
		)
		return err
	}
	for _, operation := range req.Operations {
		switch operation.Operation {
		case "add":
			if len(operation.Permissions) > 0 {

				for _, permission := range operation.Permissions {
					tuple := entity.Tuple{Object: permission.Resource, Relation: permission.Action, SubjectId: permission.Role}
					exists, err := s.CheckTuple(ctx, org.Key, tuple)
					if exists {
						s.logger.Info("Tuple is already existed in keto. Hence skipping the tuple creation.",
							zap.String("subject", tuple.SubjectId),
							zap.String("object", tuple.Object),
							zap.String("relation", tuple.Relation),
						)
						continue
					}
					if err != nil {
						return err
					}
					if err = s.CreateTuple(ctx, org.Key, tuple); err != nil {
						return err
					}
				}
			}
		case "remove":
			if len(operation.Permissions) > 0 {
				for _, permission := range operation.Permissions {
					tuple := entity.Tuple{Object: permission.Resource, Relation: permission.Action, SubjectId: permission.Role}
					exists, err := s.CheckTuple(ctx, org.Key, tuple)
					if !exists {
						s.logger.Info("Tuple is not exists in keto. Hence skipping the tuple deletion.",
							zap.String("subject", tuple.SubjectId),
							zap.String("object", tuple.Object),
							zap.String("relation", tuple.Relation),
						)
						continue
					}
					if err != nil {
						return err
					}
					if err = s.DeleteTuple(ctx, org.Key, tuple); err != nil {
						return err
					}
				}
			}
		}
	}

	// Here we are flushing the cache after every patch request. TODO: Need to find a better way to do this.
	s.permissionCache.FlushAll(ctx)

	return nil
}

type CheckActionsRequest struct {
	Role     string   `json:"role"`
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

func (s service) CheckActions(ctx context.Context, org_id string, request CheckActionsRequest) ([]string, error) {

	org, err := s.repo.GetOrganization(ctx, org_id)
	if err != nil {
		s.logger.Error("Error while getting organization from database.",
			zap.String("organization_id", org_id),
		)
		return []string{}, err
	}
	allowed_actions := []string{}
	for _, action := range request.Actions {
		tuple := entity.Tuple{Object: request.Resource, Relation: action, SubjectId: request.Role}
		bool, err := s.CheckTuple(ctx, org.Key, tuple)
		if err != nil {
			return []string{}, &util.SystemError{}
		}
		if bool {
			allowed_actions = append(allowed_actions, action)
		}
	}
	return allowed_actions, nil
}

// Get qualified tuple with organization prefix.
func qualifiedTuple(org string, tuple entity.Tuple) entity.Tuple {

	tuple.Object = org + "/" + tuple.Object
	tuple.SubjectId = org + "/" + tuple.SubjectId
	return tuple
}
