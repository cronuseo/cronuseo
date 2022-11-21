package permission

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, resource_id string, id string) (entity.Permission, error)
	Query(ctx context.Context, resource_id string) ([]entity.Permission, error)
	Create(ctx context.Context, resource_id string, permission entity.Permission) error
	Update(ctx context.Context, resource_id string, permission entity.Permission) error
	Delete(ctx context.Context, resource_id string, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Get(ctx context.Context, resource_id string, id string) (entity.Permission, error) {
	permission := entity.Permission{}
	err := r.db.Get(&permission, "SELECT * FROM permission WHERE resource_id = $1 AND permission_id = $2", resource_id, id)
	return permission, err
}

func (r repository) Create(ctx context.Context, resource_id string, permission entity.Permission) error {

	stmt, err := r.db.Prepare("INSERT INTO permission(permission_key,name,resource_id,permission_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(permission.Key, permission.Name, resource_id, permission.ID)
	if err != nil {
		return err
	}
	return nil

}

func (r repository) Update(ctx context.Context, resource_id string, permission entity.Permission) error {
	stmt, err := r.db.Prepare("UPDATE permission SET name = $1 HERE resource_id = $2 AND permission_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(permission.Name, resource_id, permission.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, resource_id string, id string) error {
	stmt, err := r.db.Prepare("DELETE FROM permission WHERE resource_id = $3 AND permission_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource_id, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Query(ctx context.Context, resource_id string) ([]entity.Permission, error) {
	resources := []entity.Permission{}
	err := r.db.Select(&resources, "SELECT * FROM permission WHERE resource_id = $1", resource_id)
	return resources, err
}
