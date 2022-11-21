package entity

type User struct {
	ID        string `json:"user_id" db:"user_id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"firstname" db:"firstname"`
	LastName  string `json:"lastname" db:"lastname"`
	OrgID     string `json:"org_id" db:"org_id"`
}
