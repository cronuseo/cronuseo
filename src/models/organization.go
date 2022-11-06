package models

// @Description Organization information
type Organization struct {
	ID   int    `json:"org_id" gorm:"primary_key"`
	Key  string `json:"key" validate:"required,min=4"`
	Name string `json:"name" validate:"required,min=4"`
}
