package check

import (
	"context"
	"encoding/json"

	"github.com/shashimalcse/cronuseo/internal/util"
	"github.com/shashimalcse/tunnel_go"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type CheckDetails struct {
	Roles          []primitive.ObjectID
	Policies       []primitive.ObjectID
	UserProperties map[string]interface{}
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
	checkDetails, err := s.repo.GetCheckDetails(ctx, org_identifier, req.Identifier)
	if err != nil {
		return false, err
	}
	allow := false
	if len(checkDetails.Roles) > 0 {
		role_permissions, err := s.repo.GetRolePermissions(ctx, org_identifier, checkDetails.Roles)
		if err != nil {
			return false, err
		}
		for _, permission := range *role_permissions {
			if permission.Resource == req.Resource && permission.Action == req.Action {
				allow = true
			}
		}
	}
	if !skipValidation {
		properties, err := json.Marshal(*&checkDetails.UserProperties)
		if err != nil {
			return false, err
		}
		active_policies, err := s.repo.GetActivePolicyVersionContents(ctx, org_identifier, checkDetails.Policies)
		for _, policy := range active_policies {
			result := tunnel_go.ValidateTunnelPolicy(policy, string(properties))
			if !result {
				return false, nil
			}
		}
	}
	return allow, nil
}
