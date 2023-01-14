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
	GetUserByUsername(ctx context.Context, username string) (entity.AdminUser, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

func (r repository) Create(ctx context.Context, adminUser entity.AdminUser) error {

	stmt, err := r.db.Prepare("INSERT INTO org_admin_user(user_id,username,password,is_super) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(adminUser.ID, adminUser.Username, adminUser.Password, adminUser.IsSuper)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Get(ctx context.Context, id string) (entity.AdminUser, error) {
	user := entity.AdminUser{}
	err := r.db.Get(&user, "SELECT * FROM org_admin_user WHERE user_id = $1", id)
	return user, err
}

func (r repository) GetUserByUsername(ctx context.Context, username string) (entity.AdminUser, error) {
	user := entity.AdminUser{}
	err := r.db.Get(&user, "SELECT * FROM org_admin_user WHERE username = $1", username)
	return user, err
}

func (r repository) ExistByUsername(ctx context.Context, username string) (bool, error) {
	exists := false
	err := r.db.QueryRow("SELECT exists (SELECT user_id FROM org_admin_user WHERE username = $1)", username).Scan(&exists)
	return exists, err
}
