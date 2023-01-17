package entity

type AdminUser struct {
	LogicalKey string `json:"-" db:"id"`
	ID         string `json:"user_id" db:"user_id"`
	Username   string `json:"username" db:"username"`
	Password   []byte `json:"-" db:"password"`
	OrgID      string `json:"org_id" db:"org_id"`
	IsSuper    bool   `json:"-" db:"is_super"`
	CreatedAt  string `json:"-" db:"created_at"`
	UpdatedAt  string `json:"-" db:"updated_at"`
}
