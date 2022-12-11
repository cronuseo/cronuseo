package entity

type Resource struct {
	LogicalKey string `json:"-" db:"id"`
	ID         string `json:"resource_id" db:"resource_id"`
	Key        string `json:"resource_key" db:"resource_key"`
	Name       string `json:"name" db:"name"`
	OrgID      string `json:"org_id" db:"org_id"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	UpdatedAt  string `json:"updated_at" db:"updated_at"`
}
