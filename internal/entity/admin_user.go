package entity

import "database/sql"

type AdminUser struct {
	LogicalKey string         `json:"-" db:"id"`
	ID         string         `json:"user_id" db:"user_id"`
	Username   string         `json:"username" db:"username"`
	Password   []byte         `json:"password" db:"password"`
	OrgID      sql.NullString `json:"org_id" db:"org_id"`
	IsSuper    bool           `json:"is_super" db:"is_super"`
	CreatedAt  string         `json:"created_at" db:"created_at"`
	UpdatedAt  string         `json:"updated_at" db:"updated_at"`
}
