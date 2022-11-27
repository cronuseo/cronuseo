package role

import (
	"context"
	"log"

	"github.com/shashimalcse/cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, org_id string, id string) (entity.Role, error)
	Query(ctx context.Context, org_id string) ([]entity.Role, error)
	Create(ctx context.Context, org_id string, role entity.Role) error
	Update(ctx context.Context, org_id string, role entity.Role) error
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

func (r repository) Get(ctx context.Context, org_id string, id string) (entity.Role, error) {
	role := entity.Role{}
	err := r.db.Get(&role, "SELECT * FROM org_role WHERE org_id = $1 AND role_id = $2", org_id, id)
	return role, err
}

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
	// add users
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

func (r repository) Delete(ctx context.Context, org_id string, id string) error {
	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}
	{
		stmt, err := tx.Prepare("DELETE FROM org_role WHERE org_id = $3 AND role_id = $1")
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
		stmt, err := tx.Prepare(`DELETE FROM org_role WHERE role_id = $1`)

		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(id)

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

func (r repository) Query(ctx context.Context, org_id string) ([]entity.Role, error) {
	roles := []entity.Role{}
	err := r.db.Select(&roles, "SELECT * FROM org_role WHERE org_id = $1", org_id)
	return roles, err
}

func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT role_id FROM org_role WHERE role_id = $1)", id).Scan(&exists)
	return exists, err
}

func (r repository) ExistByKey(ctx context.Context, key string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT role_id FROM org_role WHERE role_key = $1)", key).Scan(&exists)
	return exists, err
}
