package keto

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Service interface {
	CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
}

type Tuple struct {
	entity.Tuple
}

type service struct {
	writeClient rts.WriteServiceClient
	readClient  rts.ReadServiceClient
	checkClient rts.CheckServiceClient
}

type KetoClients struct {
	WriteClient rts.WriteServiceClient
	ReadClient  rts.ReadServiceClient
	CheckClient rts.CheckServiceClient
}

func NewService(ketoClients KetoClients) Service {
	return service{writeClient: ketoClients.WriteClient, readClient: ketoClients.ReadClient, checkClient: ketoClients.CheckClient}
}

func (s service) CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {
	_, err := s.writeClient.TransactRelationTuples(context.Background(), &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: namespace,
					Object:    tuple.Object,
					Relation:  tuple.Relation,
					Subject:   rts.NewSubjectID(tuple.SubjectId),
				},
			},
		},
	})
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
	return nil
}

func (s service) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error) {
	check, err := s.checkClient.Check(context.Background(), &rts.CheckRequest{
		Namespace: namespace,
		Object:    tuple.Object,
		Relation:  tuple.Relation,
		Subject:   rts.NewSubjectID(tuple.SubjectId),
	})
	return check.Allowed, err
}

func (s service) GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {
	//do the list operation
	return []string{}, nil
}

func (s service) GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {
	//do the list operation
	return []string{}, nil
}
