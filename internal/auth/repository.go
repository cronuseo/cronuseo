package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/shashimalcse/cronuseo/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, adminUser entity.AdminUser) error
	ExistByUsername(ctx context.Context, username string) (bool, error)
	Get(ctx context.Context, id string) (entity.AdminUser, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Create(ctx context.Context, adminUser entity.AdminUser) error {

	stmt, err := r.db.Prepare("INSERT INTO org(org_key,name,org_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Get(ctx context.Context, id string) (entity.AdminUser, error) {
	user := entity.AdminUser{}
	err := r.db.Get(&user, "SELECT * FROM admin_user WHERE user_id = $1", id)
	return user, err
}

func (r repository) ExistByUsername(ctx context.Context, username string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT user_id FROM admin_user WHERE username = $1)", username).Scan(&exists)
	return exists, err
}
