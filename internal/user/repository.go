package user

import (
	"context"
	"log"

	"github.com/shashimalcse/cronuseo/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, org_id string, id string) (entity.User, error)
	Query(ctx context.Context, org_id string, filter Filter) ([]entity.User, error)
	Create(ctx context.Context, org_id string, user entity.User) error
	Update(ctx context.Context, org_id string, user entity.User) error
	Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) error
	Delete(ctx context.Context, org_id string, id string) error
	ExistByID(ctx context.Context, id string) (bool, error)
	ExistByKey(ctx context.Context, username string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {

	return repository{db: db}
}

// Get user by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (entity.User, error) {

	user := entity.User{}
	err := r.db.Get(&user, "SELECT * FROM org_user WHERE org_id = $1 AND user_id = $2", org_id, id)
	return user, err
}

// Create new user.
func (r repository) Create(ctx context.Context, org_id string, user entity.User) error {

	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}
	{
		stmt, err := tx.Prepare("INSERT INTO org_user(username,firstname,lastname,org_id,user_id) VALUES($1, $2, $3, $4, $5)")
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(user.Username, user.FirstName, user.LastName, org_id, user.ID); err != nil {
			return err
		}
	}
	// Add roles.
	if len(user.Roles) > 0 {
		stmt, err := tx.Prepare("INSERT INTO user_role(user_id,role_id) VALUES($1, $2)")
		if err != nil {
			return err
		}
		defer stmt.Close()
		for _, role := range user.Roles {
			if _, err = stmt.Exec(user.ID, role.ID); err != nil {
				log.Fatal(err)
			}

		}
	}
	{
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil

}

// Update user.
func (r repository) Update(ctx context.Context, org_id string, user entity.User) error {

	stmt, err := r.db.Prepare("UPDATE org_user SET firstname = $1, lastname = $2 WHERE org_id = $3 AND user_id = $4")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(user.FirstName, user.LastName, org_id, user.ID); err != nil {
		return err
	}
	return nil
}

// Delete user.
func (r repository) Delete(ctx context.Context, org_id string, id string) error {

	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}

	// Delete roles assigned to the user.
	{
		stmt, err := tx.Prepare("DELETE FROM user_role WHERE user_id = $1")
		if err != nil {
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(id); err != nil {
			return err
		}
	}

	// Delete user.
	{
		stmt, err := tx.Prepare("DELETE FROM org_user WHERE org_id = $1 AND user_id = $2")
		if err != nil {
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(org_id, id); err != nil {
			return err
		}
	}

	{
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// Get all users.
func (r repository) Query(ctx context.Context, org_id string, filter Filter) ([]entity.User, error) {

	users := []entity.User{}
	name := filter.Name + "%"
	err := r.db.Select(&users, "SELECT * FROM org_user WHERE org_id = $1 AND username LIKE $2 AND id > $3 ORDER BY id LIMIT $4", org_id, name, filter.Cursor, filter.Limit)
	return users, err
}

// Patch user.
func (r repository) Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) error {
	tx, err := r.db.DB.Begin()

	if err != nil {
		return err
	}

	{
		for _, operation := range req.Operations {
			switch operation.Path {
			case "roles":
				{
					switch operation.Operation {
					case "add":
						if len(operation.Values) > 0 {
							stmt, err := tx.Prepare("INSERT INTO user_role(user_id,role_id) VALUES($1, $2)")
							if err != nil {
								return err
							}
							defer stmt.Close()
							for _, roleId := range operation.Values {
								exists, err := r.isRoleAlreadyAssigned(roleId.Value, id)
								if exists {
									continue
								}
								if err != nil {
									return err
								}
								_, err = stmt.Exec(id, roleId.Value)
								if err != nil {
									return err
								}
							}
						}
					case "remove":
						if len(operation.Values) > 0 {
							stmt, err := tx.Prepare("DELETE FROM user_role WHERE user_id = $1 AND role_id = $2")
							if err != nil {
								return err
							}
							defer stmt.Close()
							for _, roleId := range operation.Values {
								exists, err := r.isRoleAlreadyAssigned(roleId.Value, id)
								if !exists {
									continue
								}
								if err != nil {
									return err
								}
								_, err = stmt.Exec(id, roleId.Value)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
	}

	{
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// Check if role is already assigned to the user.
func (r repository) isRoleAlreadyAssigned(roleId string, userId string) (bool, error) {

	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT role_id FROM user_role WHERE role_id = $1 AND user_id = $2)", roleId, userId).Scan(&exists)
	return exists, err
}

// Get user by id.
func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {

	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT user_id FROM org_user WHERE user_id = $1)", id).Scan(&exists)
	return exists, err
}

// Get user by username.
func (r repository) ExistByKey(ctx context.Context, username string) (bool, error) {

	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT user_id FROM org_user WHERE username = $1)", username).Scan(&exists)
	return exists, err
}
