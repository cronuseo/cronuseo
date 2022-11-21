package organization

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Organization, error)
	Query(ctx context.Context) ([]entity.Organization, error)
	Create(ctx context.Context, organization entity.Organization) error
	Update(ctx context.Context, organization entity.Organization) error
	Delete(ctx context.Context, id string) error
	ExistByID(ctx context.Context, id string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Get(ctx context.Context, id string) (entity.Organization, error) {
	organization := entity.Organization{}
	err := r.db.Get(&organization, "SELECT * FROM org WHERE org_id = $1", id)
	return organization, err
}

func (r repository) Create(ctx context.Context, organization entity.Organization) error {

	stmt, err := r.db.Prepare("INSERT INTO org(org_key,name,org_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(organization.Key, organization.Name, organization.ID)
	if err != nil {
		return err
	}
	return nil

}

func (r repository) Update(ctx context.Context, organization entity.Organization) error {
	stmt, err := r.db.Prepare("UPDATE org SET name = $1 WHERE org_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(organization.Name, organization.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Delete(ctx context.Context, id string) error {
	stmt, err := r.db.Prepare("DELETE FROM org WHERE org_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Query(ctx context.Context) ([]entity.Organization, error) {
	organizations := []entity.Organization{}
	err := r.db.Select(&organizations, "SELECT * FROM org")
	return organizations, err
}

func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT org_id FROM org WHERE org_id = $1)", id).Scan(&exists)
	return exists, err
}
