package entity

type Role struct {
	ID    string   `json:"role_id" db:"role_id"`
	Key   string   `json:"role_key" db:"role_key"`
	Name  string   `json:"name" db:"name"`
	OrgID string   `json:"org_id" db:"org_id"`
	Users []UserID `json:"users"`
}

type UserID struct {
	ID string `json:"user_id" db:"user_id"`
}
