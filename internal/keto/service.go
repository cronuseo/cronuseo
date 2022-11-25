package keto

import (
	"context"
	"errors"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"google.golang.org/protobuf/internal/errors"
)

type Service interface {
	CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
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
		return errors.New("")
	}
	if err != nil {
		return errors.New("")
	}

	tuple = qualifiedTuple(org, tuple)
	return s.repo.CreateTuple(ctx, org, namespace, tuple)
}

func (s service) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error) {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.CheckTuple(ctx, org, namespace, tuple)

}
func (s service) GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.GetObjectListBySubject(ctx, org, namespace, tuple)

}

func (s service) GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.GetSubjectListByObject(ctx, org, namespace, tuple)

}
func (s service) DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {

	tuple = qualifiedTuple(org, tuple)
	return s.repo.DeleteTuple(ctx, org, namespace, tuple)
}

func qualifiedTuple(org string, tuple entity.Tuple) entity.Tuple {

	tuple.Object = org + "/" + tuple.Object
	tuple.SubjectId = org + "/" + tuple.SubjectId
	return tuple
}
