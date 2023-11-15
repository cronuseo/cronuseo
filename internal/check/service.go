package check

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"
)

type Service interface {
	Check(ctx context.Context, org_identifier string, req CheckRequest, apiKey string) (bool, error)
}

type CheckRequest struct {
	Identifier string `json:"identifier"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {

	return service{repo: repo, logger: logger}
}

func (s service) Check(ctx context.Context, org_identifier string, req CheckRequest, apiKey string) (bool, error) {

	// Check resource already exists.
	validated, _ := s.repo.ValidateAPIKey(ctx, org_identifier, apiKey)
	if !validated {
		s.logger.Debug("API_KEY is not valid.")
		return false, &util.UnauthorizedError{}
	}

	user_roles, err := s.repo.GetUserRoles(ctx, org_identifier, req.Identifier)
	if err != nil {
		s.logger.Error("Error while retrieving user roles.",
			zap.String("org_identifier", org_identifier), zap.String("identifier", req.Identifier))
		return false, err
	}
	group_roles, err := s.repo.GetGroupRoles(ctx, org_identifier, req.Identifier)
	if err != nil {
		s.logger.Error("Error while retrieving group roles.",
			zap.String("org_identifier", org_identifier), zap.String("identifier", req.Identifier))
		return false, err
	}

	user_roles_map := append(*user_roles, *group_roles...)

	role_permissions, err := s.repo.GetRolePermissions(ctx, org_identifier, user_roles_map)
	allow := false
	for _, permission := range *role_permissions {
		if permission.Resource == req.Resource && permission.Action == req.Action {
			allow = true
		}
	}
	return allow, nil
}
