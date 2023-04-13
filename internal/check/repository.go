package check

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	ValidateAPIKey(ctx context.Context, org string, apiKey string) (bool, error)
}

type repository struct {
	mongodb *mongo.Database
}

func NewRepository(mongodb *mongo.Database) Repository {
	return repository{mongodb: mongodb}
}

func (r repository) ValidateAPIKey(ctx context.Context, org string, apiKey string) (bool, error) {

	filter := bson.M{"identifier": org, "api_key": apiKey}

	// Search for the resource in the "organizations" collection
	count, err := r.mongodb.Collection("organizations").CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
