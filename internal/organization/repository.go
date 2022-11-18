package organization

import (
	"context"
	"cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Organization, error)
	Query(ctx context.Context) ([]entity.Organization, error)
	Create(ctx context.Context, organization entity.Organization) error
	Update(ctx context.Context, organization entity.Organization) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Get(ctx context.Context, id string) (entity.Organization, error) {
	var orgnization entity.Organization
	err := r.db.Get(orgnization, "SELECT * FROM org WHERE org_id = $1", id)
	return orgnization, err
}

func (r repository) Create(ctx context.Context, orgnization entity.Organization) error {

	var org_id string

	tx, err := r.db.Begin()

	if err != nil {
		return err
	}
	// add group
	{
		stmt, err := tx.Prepare(`INSERT INTO org(org_key,name) VALUES($1, $2) RETURNING org_id`)

		if err != nil {
			return err
		}

		defer stmt.Close()

		err = stmt.QueryRow(orgnization.Key, orgnization.Name).Scan(&org_id)

		if err != nil {
			return err
		}
	}

	// commit changes
	{
		err := tx.Commit()

		if err != nil {
			return err
		}
	}
	orgnization.ID = org_id
	return nil

}

// Update saves the changes to an organization in the database.
func (r repository) Update(ctx context.Context, orgnization entity.Organization) error {
	stmt, err := r.db.Prepare("UPDATE organization SET name = $1 WHERE org_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(orgnization.Name, orgnization.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes an organization with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	stmt, err := r.db.Prepare("DELETE FROM organization WHERE org_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// Query retrieves the organization records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context) ([]entity.Organization, error) {
	var orgnizations []entity.Organization
	err := r.db.Select(orgnizations, "SELECT * FROM organization")
	return orgnizations, err
}
