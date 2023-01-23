package role

import (
	"context"
	"log"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/permission"

	"github.com/jmoiron/sqlx"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type Repository interface {
	Get(ctx context.Context, org_id string, id string) (entity.Role, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]entity.Role, error)
	QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]entity.Role, error)
	Create(ctx context.Context, org_id string, role entity.Role) error
	Update(ctx context.Context, org_id string, role entity.Role) error
	Delete(ctx context.Context, role entity.Role) error
	ExistByID(ctx context.Context, id string) (bool, error)
	ExistByKey(ctx context.Context, key string) (bool, error)
}

type repository struct {
	db          *sqlx.DB
	writeClient rts.WriteServiceClient
}

func NewRepository(db *sqlx.DB, writeClient rts.WriteServiceClient) Repository {

	return repository{db: db, writeClient: writeClient}
}

// Get role by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (entity.Role, error) {

	role := entity.Role{}
	err := r.db.Get(&role, "SELECT * FROM org_role WHERE org_id = $1 AND role_id = $2", org_id, id)
	return role, err
}

// Create new role.
func (r repository) Create(ctx context.Context, org_id string, role entity.Role) error {

	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}
	{
		stmt, err := tx.Prepare("INSERT INTO org_role(role_key,name,org_id,role_id) VALUES($1, $2, $3, $4)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(role.Key, role.Name, org_id, role.ID)
		if err != nil {
			return err
		}
	}
	// Assign users to the role.
	if len(role.Users) > 0 {
		stmt, err := tx.Prepare("INSERT INTO user_role(user_id,role_id) VALUES($1, $2)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		for _, user := range role.Users {
			_, err = stmt.Exec(user.ID, role.ID)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}
	return nil

}

// Update role.
func (r repository) Update(ctx context.Context, org_id string, role entity.Role) error {

	stmt, err := r.db.Prepare("UPDATE org_role SET name = $1 HERE org_id = $2 AND role_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(role.Name, org_id, role.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete role.
func (r repository) Delete(ctx context.Context, role entity.Role) error {

	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}
	// Delete all permissions of the role.
	{
		organization := entity.Organization{}
		err := r.db.Get(&organization, "SELECT * FROM org WHERE org_id = $1", role.OrgID)
		if err != nil {
			return err
		}
		resources := []entity.Resource{}
		err = r.db.Select(&resources, "SELECT * FROM org_resource WHERE org_id = $1", role.OrgID)
		if err != nil {
			return err
		}
		for _, resource := range resources {
			actions := []entity.Action{}
			err := r.db.Select(&actions, "SELECT * FROM res_action WHERE resource_id = $1", resource.ID)
			if err != nil {
				return err
			}
			for _, action := range actions {
				tuple := entity.Tuple{SubjectId: role.Key, Relation: action.Key, Object: resource.Key}
				qTuple := qualifiedTuple(organization.Key, tuple)
				err = r.DeleteTuple(ctx, qTuple)
				if err != nil {
					return err
				}
			}
		}

	}

	// Delete all users relationship of the role.
	{
		stmt, err := tx.Prepare(`DELETE FROM user_role WHERE role_id = $1`)

		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(role.ID)

		if err != nil {
			return err
		}
	}

	// Delete the role.
	{
		stmt, err := tx.Prepare("DELETE FROM org_role WHERE org_id = $1 AND role_id = $2")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(role.OrgID, role.ID)
		if err != nil {
			return err
		}
	}

	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}
	return nil
}

// Query roles.
func (r repository) Query(ctx context.Context, org_id string, filter Filter) ([]entity.Role, error) {

	roles := []entity.Role{}
	name := filter.Name + "%"
	err := r.db.Select(&roles, "SELECT * FROM org_role WHERE org_id = $1 AND name LIKE $2 AND id > $3 ORDER BY id LIMIT $4", org_id, name, filter.Cursor, filter.Limit)
	return roles, err
}

// Query roles by user id.
func (r repository) QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]entity.Role, error) {
	roles := []entity.Role{}
	name := filter.Name + "%"
	err := r.db.Select(&roles, "SELECT org_role.id, org_role.role_id, org_role.role_key, org_role.name, org_role.org_id, org_role.created_at, org_role.updated_at FROM org_role INNER JOIN user_role ON org_role.role_id = user_role.role_id WHERE org_role.org_id = $1 AND user_role.user_id = $2 AND org_role.name LIKE $3 AND org_role.id > $4 ORDER BY org_role.id LIMIT $5", org_id, user_id, name, filter.Cursor, filter.Limit)
	return roles, err
}

// Check if role exists by id.
func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT role_id FROM org_role WHERE role_id = $1)", id).Scan(&exists)
	return exists, err
}

// Check if role exists by key.
func (r repository) ExistByKey(ctx context.Context, key string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT role_id FROM org_role WHERE role_key = $1)", key).Scan(&exists)
	return exists, err
}

func qualifiedTuple(org string, tuple entity.Tuple) entity.Tuple {

	tuple.Object = org + "/" + tuple.Object
	tuple.SubjectId = org + "/" + tuple.SubjectId
	return tuple
}

// When role is deleted, we need to delete all permissions that associated with the role.
// Here we use keto to delete the permissions.
func (r repository) DeleteTuple(ctx context.Context, tuple entity.Tuple) error {

	_, err := r.writeClient.TransactRelationTuples(ctx, &rts.TransactRelationTuplesRequest{
		RelationTupleDeltas: []*rts.RelationTupleDelta{
			{
				Action: rts.RelationTupleDelta_ACTION_DELETE,
				RelationTuple: &rts.RelationTuple{
					Namespace: permission.NAMESPACE,
					Object:    tuple.Object,
					Relation:  tuple.Relation,
					Subject:   rts.NewSubjectID(tuple.SubjectId),
				},
			},
		},
	})
	return err
}
