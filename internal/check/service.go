package check

import (
	"context"
	"strconv"
	"strings"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.uber.org/zap"
)

type Service interface {
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple, apiKey string) (bool, error)
	CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.CheckRequestWithUser, apiKey string) (bool, error)
	CheckPermissions(ctx context.Context, org string, namespace string, tuple entity.CheckRequestWithPermissions, apiKey string) ([]string, error)
	CheckAll(ctx context.Context, org string, namespace string, tuple entity.CheckRequestAll, apiKey string) (entity.CheckAllResponse, error)
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
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

	return service{repo: repo, permissionCache: cache, logger: logger}
}

// Checks permission is allowed or not.
func (s service) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple, apiKey string) (bool, error) {

	// Checking API Key is valid or not.
	api_key, _ := s.permissionCache.GetAPIKey(ctx, "API_KEY")
	if api_key == "" {
		s.logger.Debug("API_KEY is not in the cache.")
		orgObject, err := s.repo.GetOrganizationByKey(ctx, org)
		if err != nil {
			s.logger.Error("Error while getting organization from database.",
				zap.String("organization key", org),
			)
			return false, &util.NotFoundError{Path: "Organization"}
		}
		api_key = orgObject.API_KEY
		s.logger.Debug("Adding API_KEY into the cache.")
		s.permissionCache.SetAPIKey(ctx, "API_KEY", api_key)
	}

	if apiKey != api_key {
		s.logger.Debug("API_KEY is not valid.")
		return false, &util.UnauthorizedError{}
	}

	tuple = qualifiedTuple(org, tuple)
	value, _ := s.permissionCache.Get(ctx, tuple)
	if value == "true" {
		s.logger.Debug("Checking permission is in the cache.")
		return true, nil
	}
	if value == "false" {
		s.logger.Debug("Checking permission is in the cache.")
		return false, nil
	}
	s.logger.Debug("Checking permission is not in the cache. Hence checking from keto.")
	allow, err := s.repo.CheckTuple(ctx, org, namespace, tuple)
	if err != nil {
		s.logger.Error("Error while checking tuple with keto.",
			zap.String("organization", org),
			zap.String("subject", tuple.SubjectId),
			zap.String("object", tuple.Object),
			zap.String("relation", tuple.Relation),
		)
		return false, err
	}
	b := strconv.FormatBool(allow)
	s.permissionCache.Set(ctx, tuple, b)
	return allow, nil

}

// Checks if the user has the permission to perform the action by username.
func (s service) CheckByUsername(ctx context.Context, org string, namespace string, tupleUsername entity.CheckRequestWithUser, apiKey string) (bool, error) {

	// Checking API Key is valid or not.
	api_key, _ := s.permissionCache.GetAPIKey(ctx, "API_KEY")
	if api_key == "" {
		s.logger.Debug("API_KEY is not in the cache.")
		orgObject, err := s.repo.GetOrganizationByKey(ctx, org)
		if err != nil {
			s.logger.Error("Error while getting organization from database.",
				zap.String("organization key", org),
			)
			return false, &util.NotFoundError{Path: "Organization"}
		}
		api_key = orgObject.API_KEY
		s.logger.Debug("Adding API_KEY into the cache.")
		s.permissionCache.SetAPIKey(ctx, "API_KEY", api_key)
	}

	if apiKey != api_key {
		s.logger.Debug("API_KEY is not valid.")
		return false, &util.UnauthorizedError{}
	}

	tuple := entity.Tuple{SubjectId: tupleUsername.Username, Object: tupleUsername.Resource, Relation: tupleUsername.Permission}
	qTuple := qualifiedTuple(org, tuple)
	value, _ := s.permissionCache.Get(ctx, qTuple)
	if value == "true" {
		s.logger.Debug("Checking permission is in the cache.")
		return true, nil
	}
	if value == "false" {
		s.logger.Debug("Checking permission is in the cache.")
		return false, nil
	}

	s.logger.Debug("Checking permission is not in the cache. Hence checking from keto.")
	roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, qTuple)
	if err != nil {
		s.logger.Error("Error while retrieving subjects with keto.",
			zap.String("organization", org),
			zap.String("subject", tuple.SubjectId),
			zap.String("object", tuple.Object),
			zap.String("relation", tuple.Relation),
		)
		return false, err
	}
	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, tuple.SubjectId)
	if err != nil {
		s.logger.Error("Error while retrieving roles from database.",
			zap.String("organization", org),
			zap.String("username", tuple.SubjectId),
		)
		return false, err
	}
	allow := false
	for _, val := range roles_from_db {
		if contains(roles_from_keto, val) {
			s.logger.Debug("User has the role. hence we are allowing the permission.")
			allow = true
		}
	}
	b := strconv.FormatBool(allow)
	s.permissionCache.Set(ctx, qTuple, b)
	return allow, nil
}

func (s service) CheckPermissions(ctx context.Context, org string, namespace string, permissions entity.CheckRequestWithPermissions, apiKey string) ([]string, error) {

	api_key, _ := s.permissionCache.GetAPIKey(ctx, "API_KEY")
	if api_key == "" {
		orgObject, err := s.repo.GetOrganizationByKey(ctx, org)
		if err != nil {
			return []string{}, err
		}
		api_key = orgObject.API_KEY
		s.permissionCache.SetAPIKey(ctx, "API_KEY", api_key)
	}
	if apiKey != api_key {
		return []string{}, &util.UnauthorizedError{}
	}

	allowedPermissions := []string{}
	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, permissions.Username)
	if err != nil {
		return []string{}, err
	}
	for _, permission := range permissions.Permissions {
		tuple := entity.Tuple{Object: permissions.Resource, Relation: permission.Permission}
		qTuple := qualifiedTuple(org, tuple)
		roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, qTuple)
		if err != nil {
			return []string{}, err
		}
		for _, val := range roles_from_db {
			if contains(roles_from_keto, val) {
				allowedPermissions = append(allowedPermissions, permission.Permission)
			}
		}
	}
	return allowedPermissions, nil
}

func (s service) CheckAll(ctx context.Context, org string, namespace string, tuple entity.CheckRequestAll, apiKey string) (entity.CheckAllResponse, error) {

	api_key, _ := s.permissionCache.GetAPIKey(ctx, "API_KEY")
	if api_key == "" {
		orgObject, err := s.repo.GetOrganizationByKey(ctx, org)
		if err != nil {
			return entity.CheckAllResponse{}, err
		}
		api_key = orgObject.API_KEY
		s.permissionCache.SetAPIKey(ctx, "API_KEY", api_key)
	}

	if apiKey != api_key {
		return entity.CheckAllResponse{}, &util.UnauthorizedError{}
	}
	response := entity.CheckAllResponse{}

	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, tuple.Username)
	if err != nil {
		return response, err
	}
	for _, resource := range tuple.Resources {
		res_name := resource.Resource
		allowedPermissions := []string{}
		for _, permission := range resource.Permissions {
			tuple := entity.Tuple{Object: res_name, Relation: permission.Permission}
			qTuple := qualifiedTuple(org, tuple)
			roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, qTuple)
			if err != nil {
				return response, err
			}
			for _, val := range roles_from_db {
				if contains(roles_from_keto, val) {
					allowedPermissions = append(allowedPermissions, permission.Permission)
				}
			}
		}
		response.Values = append(response.Values, entity.CheckAllResult{Resource: res_name, Permissions: allowedPermissions})

	}
	return response, nil
}

func (s service) GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	objects, err := s.repo.GetObjectListBySubject(ctx, org, namespace, tuple)
	if err != nil {
		return []string{}, err
	}
	values := []string{}
	for _, val := range objects {
		slc := strings.Split(val, "/")
		values = append(values, strings.TrimSpace(slc[1]))
	}
	return values, nil

}

func (s service) GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	subjects, err := s.repo.GetSubjectListByObject(ctx, org, namespace, tuple)
	if err != nil {
		return []string{}, err
	}
	values := []string{}
	for _, val := range subjects {
		slc := strings.Split(val, "/")
		values = append(values, strings.TrimSpace(slc[1]))
	}
	return values, nil

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
