package keto

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
)

type Service interface {
	CreateTuple(ctx context.Context, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, namespace string, tuple entity.Tuple) (bool, error)
	GetObjectListBySubject(ctx context.Context, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, namespace string, tuple entity.Tuple) ([]string, error)
}

type Tuple struct {
	entity.Tuple
}

type service struct {
	writeClient acl.WriteServiceClient
	readClient  acl.ReadServiceClient
	checkClient acl.CheckServiceClient
}

type KetoClients struct {
	WriteClient acl.WriteServiceClient
	ReadClient  acl.ReadServiceClient
	CheckClient acl.CheckServiceClient
}

func NewService(ketoClients KetoClients) Service {
	return service{writeClient: ketoClients.WriteClient, readClient: ketoClients.ReadClient, checkClient: ketoClients.CheckClient}
}

func (s service) CreateTuple(ctx context.Context, namespace string, tuple entity.Tuple) error {
	println("hi")
	//do the create operation
	return nil
}

func (s service) CheckTuple(ctx context.Context, namespace string, tuple entity.Tuple) (bool, error) {
	//do the check operation
	return false, nil
}

func (s service) GetObjectListBySubject(ctx context.Context, namespace string, tuple entity.Tuple) ([]string, error) {
	//do the list operation
	return []string{}, nil
}

func (s service) GetSubjectListByObject(ctx context.Context, namespace string, tuple entity.Tuple) ([]string, error) {
	//do the list operation
	return []string{}, nil
}
