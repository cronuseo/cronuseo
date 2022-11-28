package keto

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/cronuseo/internal/entity"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Repository interface {
	CreateTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	DeleteTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) error
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	GetRolesByUsername(ctx context.Context, org string, username string) ([]string, error)
}

type repo struct {
	writeClient rts.WriteServiceClient
	readClient  rts.ReadServiceClient
	checkClient rts.CheckServiceClient
	db          *sqlx.DB
}

type KetoClients struct {
	WriteClient rts.WriteServiceClient
	ReadClient  rts.ReadServiceClient
	CheckClient rts.CheckServiceClient
}

func NewRepository(ketoClients KetoClients, db *sqlx.DB) Service {
	return repo{writeClient: ketoClients.WriteClient, readClient: ketoClients.ReadClient, checkClient: ketoClients.CheckClient, db: db}
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
		obejcts = append(obejcts, rt.Subject.Ref.(*rts.Subject_Id).Id)
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

func (r repo) GetRolesByUsername(ctx context.Context, org string, username string) ([]string, error) {
	role := []string{}
	err := r.db.Get(&role, "select role_key from org_role where role_id in (select role_id from user_role where user_id in (select user_id from org_user inner join org on org_user.org_id = org.org_id where org_user.username = $1 AND org.org_key = $2));", username, org)
	return role, err
}

func (r repo) CheckByUsername(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error) {

	return false, nil
}
