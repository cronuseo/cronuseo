package check

import (
	"context"
	"fmt"
	"strconv"

	"github.com/open-policy-agent/opa/rego"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"
)

type Service interface {
	Check(ctx context.Context, org_identifier string, req CheckRequest, apiKey string) (bool, error)
}

type CheckRequest struct {
	Username string `json:"username"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type data struct {
	user_roles       []string
	role_permissions map[string][]permission
}

type permission struct {
	action   string
	resource string
}

type OPAInput struct {
	user     string
	action   string
	resource string
	data     data
}

type service struct {
	repo   Repository
	logger *zap.Logger
	query  rego.PreparedEvalQuery
}

func NewService(repo Repository, logger *zap.Logger, query rego.PreparedEvalQuery) Service {

	return service{repo: repo, logger: logger, query: query}
}

func (s service) Check(ctx context.Context, org_identifier string, req CheckRequest, apiKey string) (bool, error) {

	// Check resource already exists.
	validated, _ := s.repo.ValidateAPIKey(ctx, org_identifier, apiKey)
	if !validated {
		s.logger.Debug("API_KEY is not valid.")
		return false, &util.UnauthorizedError{}
	}

	user_roles, err := s.repo.GetUserRoles(ctx, org_identifier, req.Username)
	if err != nil {
		s.logger.Error("Error while retrieving user roles.",
			zap.String("org_identifier", org_identifier), zap.String("username", req.Username))
		return false, err
	}
	group_roles, err := s.repo.GetGroupRoles(ctx, org_identifier, req.Username)
	if err != nil {
		s.logger.Error("Error while retrieving group roles.",
			zap.String("org_identifier", org_identifier), zap.String("username", req.Username))
		return false, err
	}

	var user_roles_map []string
	for _, role := range *user_roles {
		user_roles_map = append(user_roles_map, role.Hex())
	}
	for _, role := range *group_roles {
		user_roles_map = append(user_roles_map, role.Hex())
	}

	// role_permissions, err := s.repo.GetRolePermissions(ctx, org_identifier)
	// if err != nil {
	// 	s.logger.Error("Error while retrieving user roles.",
	// 		zap.String("org_identifier", org_identifier), zap.String("username", req.Username))
	// 	return false, err
	// }
	permissions := map[string]interface{}{}
	// for _, role_permission := range *role_permissions {
	// 	var p []map[string]interface{}
	// 	for _, permission_obj := range role_permission.Permissions {
	// 		p = append(p, map[string]interface{}{"action": permission_obj.Action, "resource": permission_obj.Resource})
	// 	}
	// 	permissions[role_permission.RoleID.Hex()] = p
	// }

	opa := map[string]interface{}{
		"user":     req.Username,
		"action":   req.Action,
		"resource": req.Resource,
		"data": map[string]interface{}{
			"user_roles":       user_roles_map,
			"role_permissions": permissions,
		},
	}

	rs, err := s.query.Eval(ctx, rego.EvalInput(opa))

	if err != nil {
		return false, err
	}

	allow := fmt.Sprint(rs[0].Bindings["x"])
	boolValue, err := strconv.ParseBool(allow)
	if err != nil {
		return false, err
	}
	return boolValue, nil
}
