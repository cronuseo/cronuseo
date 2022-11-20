package entity

type User struct {
	ID        string `json:"user_id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"first_name" db:"firstname"`
	LastName  string `json:"last_name" db:"lastname"`
	OrgID     string `json:"org_id" db:"org_id"`
}
