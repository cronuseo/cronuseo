package models

type ResourceRole struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Key         string `json:"key" binding:"required,min=4"`
	Name        string `json:"name" binding:"required,min=4"`
	Description string `json:"description"`
	ResourceID  int    `json:"-" gorm:"foreignKey:ID"`
}

type ResourceRoleToGroup struct {
	GroupID        int `gorm:"foreignKey:ID"`
	ResourceRoleID int `gorm:"foreignKey:ID"`
}

type ResourceRoleToUser struct {
	UserID         int `gorm:"foreignKey:ID"`
	ResourceRoleID int `gorm:"foreignKey:ID"`
}

type ResourceRoleWithGroupsUsers struct {
	ID         int               `json:"id"`
	Key        string            `json:"key"`
	Name       string            `json:"name"`
	ResourceID int               `json:"res_id"`
	Users      []UserOnlyWithID  `json:"users"`
	Groups     []GroupOnlyWithID `json:"groups"`
}
