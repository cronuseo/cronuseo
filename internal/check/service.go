package check

import (
	"context"
	"strconv"
	"strings"

	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

type Service interface {
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple, apiKey string) (bool, error)
	CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.Tuple, apiKey string) (bool, error)
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
}

func NewService(repo Repository, cache cache.PermissionCache) Service {
	return service{repo: repo, permissionCache: cache}
}

func (s service) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple, apiKey string) (bool, error) {

	api_key, _ := s.permissionCache.GetAPIKey(ctx, "API_KEY")
	if api_key == "" {
		orgObject, err := s.repo.GetOrganizationByKey(ctx, org)
		if err != nil {
			return false, err
		}
		api_key = orgObject.API_KEY
		s.permissionCache.SetAPIKey(ctx, "API_KEY", api_key)
	}

	if apiKey != api_key {
		return false, &util.UnauthorizedError{}
	}

	tuple = qualifiedTuple(org, tuple)
	value, _ := s.permissionCache.Get(ctx, tuple)
	if value == "true" {
		return true, nil
	}
	if value == "false" {
		return false, nil
	}

	allow, err := s.repo.CheckTuple(ctx, org, namespace, tuple)
	if err != nil {
		return false, err
	}
	b := strconv.FormatBool(allow)
	s.permissionCache.Set(ctx, tuple, b)
	return allow, nil

}

func (s service) CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.Tuple, apiKey string) (bool, error) {

	api_key, _ := s.permissionCache.GetAPIKey(ctx, "API_KEY")
	if api_key == "" {
		orgObject, err := s.repo.GetOrganizationByKey(ctx, org)
		if err != nil {
			return false, err
		}
		api_key = orgObject.API_KEY
		s.permissionCache.SetAPIKey(ctx, "API_KEY", api_key)
	}
	if apiKey != api_key {
		return false, &util.UnauthorizedError{}
	}
	qTuple := qualifiedTuple(org, tuple)
	value, _ := s.permissionCache.Get(ctx, qTuple)
	if value == "true" {
		return true, nil
	}
	if value == "false" {
		return false, nil
	}

	roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, qTuple)
	if err != nil {
		return false, err
	}
	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, tuple.SubjectId)
	if err != nil {
		return false, err
	}
	allow := false
	for _, val := range roles_from_db {
		if contains(roles_from_keto, val) {
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
	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, permissions.SubjectId)
	if err != nil {
		return []string{}, err
	}
	for _, permission := range permissions.Relations {
		tuple := entity.Tuple{Object: permissions.Object, Relation: permission.Relation}
		roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, tuple)
		if err != nil {
			return []string{}, err
		}
		for _, val := range roles_from_db {
			if contains(roles_from_keto, val) {
				allowedPermissions = append(allowedPermissions, permission.Relation)
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

	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, tuple.SubjectId)
	if err != nil {
		return response, err
	}
	for _, resource := range tuple.Objects {
		res_name := resource.Object
		allowedPermissions := []string{}
		for _, permission := range resource.Relations {
			tuple := entity.Tuple{Object: res_name, Relation: permission.Relation}
			roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, tuple)
			if err != nil {
				return response, err
			}
			for _, val := range roles_from_db {
				if contains(roles_from_keto, val) {
					allowedPermissions = append(allowedPermissions, permission.Relation)
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