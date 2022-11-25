package keto

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Repository interface {
	CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
}

type repo struct {
	writeClient rts.WriteServiceClient
	readClient  rts.ReadServiceClient
	checkClient rts.CheckServiceClient
}

type KetoClients struct {
	WriteClient rts.WriteServiceClient
	ReadClient  rts.ReadServiceClient
	CheckClient rts.CheckServiceClient
}

func NewRepository(ketoClients KetoClients) Service {
	return repo{writeClient: ketoClients.WriteClient, readClient: ketoClients.ReadClient, checkClient: ketoClients.CheckClient}
}

func (r repo) CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {

	_, err := r.writeClient.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
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

func (r repo) CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error) {

	check, err := r.checkClient.Check(ctx, &rts.CheckRequest{
		Namespace: namespace,
		Object:    tuple.Object,
		Relation:  tuple.Relation,
		Subject:   rts.NewSubjectID(tuple.SubjectId),
	})
	return check.Allowed, err
}

func (r repo) GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	res, err := r.readClient.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		Query: &rts.ListRelationTuplesRequest_Query{
			Namespace: namespace,
			Relation:  tuple.Relation,
			Subject:   rts.NewSubjectID(tuple.SubjectId),
		},
	})
	if err != nil {
		return []string{}, err
	}
	obejcts := []string{}
	for _, rt := range res.RelationTuples {
		obejcts = append(obejcts, rt.Object)
	}
	return obejcts, nil
}

func (r repo) GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error) {

	res, err := r.readClient.ListRelationTuples(ctx, &rts.ListRelationTuplesRequest{
		Query: &rts.ListRelationTuplesRequest_Query{
			Namespace: namespace,
			Object:    tuple.Object,
			Relation:  tuple.Relation,
		},
	})
	if err != nil {
		return []string{}, err
	}
	obejcts := []string{}
	for _, rt := range res.RelationTuples {
		obejcts = append(obejcts, rt.GetSubject().String())
	}
	return obejcts, nil
}

func (r repo) DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error {

	_, err := r.writeClient.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
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
