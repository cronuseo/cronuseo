package mongo_entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier      string             `json:"identifier" bson:"identifier"`
	DisplayName     string             `json:"display_name" bson:"display_name"`
	API_KEY         string             `json:"api_key" bson:"api_key"`
	Resources       []Resource         `json:"resources,omitempty" bson:"resources"`
	Users           []User             `json:"users,omitempty" bson:"users"`
	Roles           []Role             `json:"roles,omitempty" bson:"roles"`
	Groups          []Group            `json:"groups,omitempty" bson:"groups"`
	RolePermissions []RolePermission   `json:"role_permissions,omitempty" bson:"role_permissions"`
}

type Resource struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	Actions     []Action           `json:"actions,omitempty" bson:"actions"`
}

type Action struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
}

type User struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Username  string               `json:"username" bson:"username"`
	Email     string               `json:"email" bson:"email"`
	FirstName string               `json:"first_name" bson:"first_name"`
	LastName  string               `json:"last_name" bson:"last_name"`
	Roles     []primitive.ObjectID `json:"roles,omitempty" bson:"roles"`
	Groups    []primitive.ObjectID `json:"groups,omitempty" bson:"groups"`
}

type Role struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Identifier  string               `json:"identifier" bson:"identifier"`
	DisplayName string               `json:"display_name" bson:"display_name"`
	Users       []primitive.ObjectID `json:"users,omitempty" bson:"users"`
	Groups      []primitive.ObjectID `json:"groups,omitempty" bson:"groups"`
}

type Group struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Identifier  string               `json:"identifier" bson:"identifier"`
	DisplayName string               `json:"display_name" bson:"display_name"`
	Users       []primitive.ObjectID `json:"users,omitempty" bson:"users"`
	Roles       []primitive.ObjectID `json:"roles,omitempty" bson:"roles"`
}

type RolePermission struct {
	RoleID      primitive.ObjectID `json:"role_id" bson:"role_id"`
	Permissions []Permission       `json:"permissions,omitempty" bson:"permissions"`
}

type Permission struct {
	Action   string `json:"action" bson:"action"`
	Resource string `json:"resource" bson:"resource"`
}
