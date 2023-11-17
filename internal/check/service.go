package check

import (
	"context"
	"encoding/json"

	"github.com/shashimalcse/cronuseo/internal/util"
	"github.com/shashimalcse/tunnel_go"
	"go.uber.org/zap"
)

type Service interface {
	Check(ctx context.Context, org_identifier string, req CheckRequest, apiKey string, skipValidation bool) (bool, error)
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

func (s service) Check(ctx context.Context, org_identifier string, req CheckRequest, apiKey string, skipValidation bool) (bool, error) {

	// Check resource already exists.
	if !skipValidation {
		validated, _ := s.repo.ValidateAPIKey(ctx, org_identifier, apiKey)
		if !validated {
			s.logger.Debug("API_KEY is not valid.")
			return false, &util.UnauthorizedError{}
		}
	}

	user_roles, err := s.repo.GetUserRoles(ctx, org_identifier, req.Identifier)
	if err != nil {
		if notFoundErr, ok := err.(*util.NotFoundError); ok {
			return false, notFoundErr
		}
		s.logger.Error("Error while retrieving user roles.",
			zap.String("org_identifier", org_identifier), zap.String("identifier", req.Identifier))
		return false, err
	}
	group_roles, err := s.repo.GetGroupRoles(ctx, org_identifier, req.Identifier)
	if err != nil {
		if notFoundErr, ok := err.(*util.NotFoundError); ok {
			return false, notFoundErr
		}
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
	if !skipValidation {
		user_properties, err := s.repo.GetUserProperties(ctx, org_identifier, req.Identifier)
		if err != nil {
			return false, err
		}
		properties, err := json.Marshal(*user_properties)
		if err != nil {
			return false, err
		}
		user_policies, err := s.repo.GetUserPolicies(ctx, org_identifier, req.Identifier)
		if err != nil {
			return false, err
		}
		user_groups, err := s.repo.GetUserGroups(ctx, org_identifier, req.Identifier)
		group_policies, err := s.repo.GetGroupPolicies(ctx, org_identifier, *user_groups)
		polcies := append(*user_policies, *group_policies...)
		active_policies, err := s.repo.GetActivePolicyVersionContents(ctx, org_identifier, polcies)
		for _, policy := range active_policies {
			result := tunnel_go.ValidateTunnelPolicy(policy, string(properties))
			if !result {
				return false, nil
			}
		}
	}
	return allow, nil
}
