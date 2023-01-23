package resource

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, org_id string, id string) (entity.Resource, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]entity.Resource, error)
	Create(ctx context.Context, org_id string, resource entity.Resource) error
	Update(ctx context.Context, org_id string, resource entity.Resource) error
	Delete(ctx context.Context, org_id string, id string) error
	ExistByID(ctx context.Context, id string) (bool, error)
	ExistByKey(ctx context.Context, key string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {

	return repository{db: db}
}

// Get resource by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (entity.Resource, error) {

	resource := entity.Resource{}
	err := r.db.Get(&resource, "SELECT * FROM org_resource WHERE org_id = $1 AND resource_id = $2", org_id, id)
	return resource, err
}

// Create new resource.
func (r repository) Create(ctx context.Context, org_id string, resource entity.Resource) error {

	stmt, err := r.db.Prepare("INSERT INTO org_resource(resource_key,name,org_id,resource_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource.Key, resource.Name, org_id, resource.ID)
	if err != nil {
		return err
	}
	return nil

}

// Update resource.
func (r repository) Update(ctx context.Context, org_id string, resource entity.Resource) error {

	stmt, err := r.db.Prepare("UPDATE org_resource SET name = $1 HERE org_id = $2 AND resource_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource.Name, org_id, resource.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete resource.
func (r repository) Delete(ctx context.Context, org_id string, id string) error {

	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}

	// Delete actions assigned to the resource
	{
		stmt, err := tx.Prepare("DELETE FROM res_action WHERE resource_id = $1")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(id)
		if err != nil {
			return err
		}
	}

	// Delete resource
	{
		stmt, err := tx.Prepare("DELETE FROM org_resource WHERE org_id = $1 AND resource_id = $2")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(org_id, id)
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

// Get all resources.
func (r repository) Query(ctx context.Context, org_id string, filter Filter) ([]entity.Resource, error) {

	resources := []entity.Resource{}
	name := filter.Name + "%"
	err := r.db.Select(&resources, "SELECT * FROM org_resource WHERE org_id = $1 AND name LIKE $2 AND id > $3 ORDER BY id LIMIT $4", org_id, name, filter.Cursor, filter.Limit)
	return resources, err
}

// Check if resource exists by id.
func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {

	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT resource_id FROM org_resource WHERE resource_id = $1)", id).Scan(&exists)
	return exists, err
}

// Check if resource exists by key.
func (r repository) ExistByKey(ctx context.Context, key string) (bool, error) {

	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT resource_id FROM org_resource WHERE resource_key = $1)", key).Scan(&exists)
	return exists, err
}
