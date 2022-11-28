package keto

import (
	"context"
	"strings"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

type Service interface {
	CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetRolesByUsername(ctx context.Context, org string, username string) ([]string, error)
	CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
}

type Tuple struct {
	entity.Tuple
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
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

func (s service) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error) {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.CheckTuple(ctx, org, namespace, tuple)

}

func (s service) CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error) {

	username := tuple.SubjectId
	tuple = qualifiedTuple(org, tuple)
	roles_from_keto, err := s.GetSubjectListByObject(ctx, org, namespace, tuple)
	if err != nil {
		return false, err
	}
	roles_from_db, err := s.repo.GetRolesByUsername(ctx, org, username)
	if err != nil {
		return false, err
	}
	for _, val := range roles_from_db {
		if contains(roles_from_keto, val) {
			return true, nil
		}
	}
	return false, nil
}

func (s service) GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	tuple = qualifiedTuple(org, tuple)
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

	tuple = qualifiedTuple(org, tuple)
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
func (s service) DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.DeleteTuple(ctx, org, namespace, tuple)
}

func (s service) GetRolesByUsername(ctx context.Context, org string, username string) ([]string, error) {

	return []string{}, nil
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
