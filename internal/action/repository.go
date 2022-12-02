package action

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, resource_id string, id string) (entity.Action, error)
	Query(ctx context.Context, resource_id string) ([]entity.Action, error)
	Create(ctx context.Context, resource_id string, action entity.Action) error
	Update(ctx context.Context, resource_id string, action entity.Action) error
	Delete(ctx context.Context, resource_id string, id string) error
	ExistByID(ctx context.Context, id string) (bool, error)
	ExistByKey(ctx context.Context, resource_id string, key string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Get(ctx context.Context, resource_id string, id string) (entity.Action, error) {
	action := entity.Action{}
	err := r.db.Get(&action, "SELECT * FROM res_action WHERE resource_id = $1 AND action_id = $2", resource_id, id)
	return action, err
}

func (r repository) Create(ctx context.Context, resource_id string, action entity.Action) error {

	stmt, err := r.db.Prepare("INSERT INTO res_action(action_key,name,resource_id,action_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(action.Key, action.Name, resource_id, action.ID)
	if err != nil {
		return err
	}
	return nil

}

func (r repository) Update(ctx context.Context, resource_id string, action entity.Action) error {
	stmt, err := r.db.Prepare("UPDATE res_action SET name = $1 WHERE resource_id = $2 AND action_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(action.Name, resource_id, action.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, resource_id string, id string) error {
	stmt, err := r.db.Prepare("DELETE FROM res_action WHERE resource_id = $3 AND action_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(resource_id, id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Query(ctx context.Context, resource_id string) ([]entity.Action, error) {
	resources := []entity.Action{}
	err := r.db.Select(&resources, "SELECT * FROM res_action WHERE resource_id = $1", resource_id)
	return resources, err
}

func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT action_id FROM res_action WHERE action_id = $1)", id).Scan(&exists)
	return exists, err
}

func (r repository) ExistByKey(ctx context.Context, resource_id string, key string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT action_id FROM res_action WHERE resource_id = $1 AND action_key = $2)", resource_id, key).Scan(&exists)
	return exists, err
}
