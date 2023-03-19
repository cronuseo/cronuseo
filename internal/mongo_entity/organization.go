package mongo_entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	API_KEY     string             `json:"api_key" bson:"api_key"`
	Resources   []Resource         `json:"resources" bson:"resources"`
}

type Resource struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	Actions     []Action           `json:"actions" bson:"actions"`
}

type Action struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"display_name" bson:"display_name"`
}
