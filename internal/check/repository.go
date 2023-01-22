package check

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Repository interface {
	CheckTuple(ctx context.Context, org string, namespace string, tuple entity.Tuple) (bool, error)
	GetObjectListBySubject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetSubjectListByObject(ctx context.Context, org string, namespace string, tuple entity.Tuple) ([]string, error)
	GetRolesByUsername(ctx context.Context, org string, username string) ([]string, error)
	GetOrganization(ctx context.Context, id string) (entity.Organization, error)
	GetOrganizationByKey(ctx context.Context, key string) (entity.Organization, error)
}

type repo struct {
	writeClient rts.WriteServiceClient
	readClient  rts.ReadServiceClient
	checkClient rts.CheckServiceClient
	db          *sqlx.DB
}

func NewRepository(ketoClients util.KetoClients, db *sqlx.DB) Repository {
	return repo{writeClient: ketoClients.WriteClient, readClient: ketoClients.ReadClient, checkClient: ketoClients.CheckClient, db: db}
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

func (r repo) GetRolesByUsername(ctx context.Context, org string, username string) ([]string, error) {
	var roles entity.Roles
	err := r.db.Select(&roles, "select * from org_role where role_id in (select role_id from user_role where user_id in (select user_id from org_user inner join org on org_user.org_id = org.org_id where org_user.username = $1 AND org.org_key = $2))", username, org)
	return roles.RoleKeys(), err
}

func (r repo) GetOrganization(ctx context.Context, id string) (entity.Organization, error) {
	organization := entity.Organization{}
	err := r.db.Get(&organization, "SELECT * FROM org WHERE org_id = $1", id)
	return organization, err
}

func (r repo) GetOrganizationByKey(ctx context.Context, key string) (entity.Organization, error) {
	organization := entity.Organization{}
	err := r.db.Get(&organization, "SELECT * FROM org WHERE org_key = $1", key)
	return organization, err
}
