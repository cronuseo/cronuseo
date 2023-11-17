package mongo_entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResourceType string

const (
	SystemResource   ResourceType = "system"
	BusinessResource ResourceType = "business"
)

type Organization struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	API_KEY     string             `json:"api_key" bson:"api_key"`
	Resources   []Resource         `json:"resources,omitempty" bson:"resources"`
	Users       []User             `json:"users,omitempty" bson:"users"`
	Roles       []Role             `json:"roles,omitempty" bson:"roles"`
	Groups      []Group            `json:"groups,omitempty" bson:"groups"`
	Polices     []Policy           `json:"policies,omitempty" bson:"policies"`
}

type Resource struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	Type        ResourceType       `json:"type" bson:"type"`
	Actions     []Action           `json:"actions,omitempty" bson:"actions"`
}

type Action struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
}

type User struct {
	ID             primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Username       string                 `json:"username" bson:"username"`
	Identifier     string                 `json:"identifier" bson:"identifier"`
	UserProperties map[string]interface{} `json:"user_properties" bson:"user_properties"`
	Roles          []primitive.ObjectID   `json:"roles,omitempty" bson:"roles"`
	Groups         []primitive.ObjectID   `json:"groups,omitempty" bson:"groups"`
	Policies       []primitive.ObjectID   `json:"policies,omitempty" bson:"policies"`
}

type Role struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Identifier  string               `json:"identifier" bson:"identifier"`
	DisplayName string               `json:"display_name" bson:"display_name"`
	Users       []primitive.ObjectID `json:"users,omitempty" bson:"users"`
	Groups      []primitive.ObjectID `json:"groups,omitempty" bson:"groups"`
	Permissions []Permission         `json:"permissions,omitempty" bson:"permissions"`
}

type Group struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Identifier  string               `json:"identifier" bson:"identifier"`
	DisplayName string               `json:"display_name" bson:"display_name"`
	Users       []primitive.ObjectID `json:"users,omitempty" bson:"users"`
	Roles       []primitive.ObjectID `json:"roles,omitempty" bson:"roles"`
	Policies    []primitive.ObjectID `json:"policies,omitempty" bson:"policies"`
}

type Permission struct {
	Action   string `json:"action" bson:"action"`
	Resource string `json:"resource" bson:"resource"`
}

type Policy struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier     string             `json:"identifier" bson:"identifier"`
	DisplayName    string             `json:"display_name" bson:"display_name"`
	ActiveVersion  string             `json:"active_version" bson:"active_version"`
	PolicyContents []PolicyContent    `json:"policy_contents" bson:"policy_contents"`
}

type PolicyContent struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Version string             `json:"version" bson:"version"`
	Policy  string             `json:"policy" bson:"policy"`
}
