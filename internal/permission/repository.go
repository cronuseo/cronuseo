package permission

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

var (
	NAMESPACE = "permission"
)

// Permission repository handle all keto and database operations.
type Repository interface {
	CreateTuple(ctx context.Context, org string, tuple entity.Tuple) error
	CheckTuple(ctx context.Context, org string, tuple entity.Tuple) (bool, error)
	DeleteTuple(ctx context.Context, org string, tuple entity.Tuple) error
	GetOrganization(ctx context.Context, id string) (entity.Organization, error)
}

type repo struct {
	writeClient rts.WriteServiceClient
	checkClient rts.CheckServiceClient
	db          *sqlx.DB
}

func NewRepository(ketoClients util.KetoClients, db *sqlx.DB) Repository {

	return repo{
		writeClient: ketoClients.WriteClient,
		checkClient: ketoClients.CheckClient,
		db:          db,
	}
}

// Create tuple in keto.
func (r repo) CreateTuple(ctx context.Context, org string, tuple entity.Tuple) error {

	_, err := r.writeClient.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_INSERT,
				RelationTuple: &rts.RelationTuple{
					Namespace: NAMESPACE,
					Object:    tuple.Object,
					Relation:  tuple.Relation,
					Subject:   rts.NewSubjectID(tuple.SubjectId),
				},
			},
		},
	})
	return err
}

// Check tuple in keto.
func (r repo) CheckTuple(ctx context.Context, org string, tuple entity.Tuple) (bool, error) {

	check, err := r.checkClient.Check(ctx, &rts.CheckRequest{
		Namespace: NAMESPACE,
		Object:    tuple.Object,
		Relation:  tuple.Relation,
		Subject:   rts.NewSubjectID(tuple.SubjectId),
	})
	return check.Allowed, err
}

// Delete tuple in keto.
func (r repo) DeleteTuple(ctx context.Context, org string, tuple entity.Tuple) error {

	_, err := r.writeClient.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: NAMESPACE,
					Object:    tuple.Object,
					Relation:  tuple.Relation,
					Subject:   rts.NewSubjectID(tuple.SubjectId),
				},
			},
		},
	})
	return err
}

// Get organization from database.
func (r repo) GetOrganization(ctx context.Context, id string) (entity.Organization, error) {

	organization := entity.Organization{}
	err := r.db.Get(&organization, "SELECT * FROM org WHERE org_id = $1", id)
	return organization, err
}
