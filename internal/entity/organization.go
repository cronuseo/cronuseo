package entity

type Organization struct {
	ID   string `json:"org_id" db:"org_id"`
	Key  string `json:"org_key" db:"org_key"`
	Name string `json:"name" db:"name"`
}
