package mongo_entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	ID          primitive.ObjectID `json:"org_id" bson:"_id,omitempty"`
	Identifier  string             `json:"identifier" bson:"identifier"`
	DisplayName string             `json:"name" bson:"name"`
	API_KEY     string             `json:"api_key" bson:"api_key"`
	// Users     []User             `bson:"users"`
	// Roles     []Role             `bson:"roles"`
	// Resources []Resource         `bson:"resources"`
}
