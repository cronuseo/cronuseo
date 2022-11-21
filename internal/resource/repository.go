package resource

import (
	"context"
	"cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, org_id string, id string) (entity.Resource, error)
	Query(ctx context.Context, org_id string) ([]entity.Resource, error)
	Create(ctx context.Context, org_id string, resource entity.Resource) error
	Update(ctx context.Context, org_id string, resource entity.Resource) error
	Delete(ctx context.Context, org_id string, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Get(ctx context.Context, org_id string, id string) (entity.Resource, error) {
	resource := entity.Resource{}
	err := r.db.Get(&resource, "SELECT * FROM org_resource WHERE org_id = $1 AND resource_id = $2", org_id, id)
	return resource, err
}

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

func (r repository) Delete(ctx context.Context, org_id string, id string) error {
	stmt, err := r.db.Prepare("DELETE FROM org_resource WHERE org_id = $3 AND resource_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(org_id, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Query(ctx context.Context, org_id string) ([]entity.Resource, error) {
	resources := []entity.Resource{}
	err := r.db.Select(&resources, "SELECT * FROM org_resource WHERE org_id = $1", org_id)
	return resources, err
}
