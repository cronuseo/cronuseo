package entity

type Role struct {
	LogicalKey string   `json:"-" db:"id"`
	ID         string   `json:"role_id" db:"role_id"`
	Key        string   `json:"role_key" db:"role_key"`
	Name       string   `json:"name" db:"name"`
	OrgID      string   `json:"org_id" db:"org_id"`
	Users      []UserID `json:"users"`
	CreatedAt  string   `json:"created_at" db:"created_at"`
	UpdatedAt  string   `json:"updated_at" db:"updated_at"`
}

type UserID struct {
	ID string `json:"user_id" db:"user_id"`
}

type Roles []Role

func (r Roles) RoleKeys() []string {
	var list []string
	for _, role := range r {
		list = append(list, role.Key)
	}
	return list
}
