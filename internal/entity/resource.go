package entity

type Resource struct {
	ID    string `json:"resource_id" db:"resource_id"`
	Key   string `json:"resource_key" db:"resource_key"`
	Name  string `json:"name" db:"name"`
	OrgID string `json:"org_id" db:"org_id"`
}
