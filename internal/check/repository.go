package check

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	ValidateAPIKey(ctx context.Context, org_identifier string, apiKey string) (bool, error)
	GetUserRoles(ctx context.Context, org_identifier string, username string) (*[]primitive.ObjectID, error)
	GetRolePermissions(ctx context.Context, org_identifier string) (*[]mongo_entity.RolePermission, error)
}

type repository struct {
	mongodb *mongo.Database
}

func NewRepository(mongodb *mongo.Database) Repository {
	return repository{mongodb: mongodb}
}

func (r repository) ValidateAPIKey(ctx context.Context, org_identifier string, apiKey string) (bool, error) {

	filter := bson.M{"identifier": org_identifier, "api_key": apiKey}

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

func (r repository) GetUserRoles(ctx context.Context, org_identifier string, username string) (*[]primitive.ObjectID, error) {

	// Define filter to find the user by its ID
	filter := bson.M{"identifier": org_identifier, "users.username": username}
	projection := bson.M{"users.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "User"}
		}
		return nil, err
	}

	return &org.Users[0].Roles, nil
}

func (r repository) GetRolePermissions(ctx context.Context, org_identifier string) (*[]mongo_entity.RolePermission, error) {

	// Define filter to find the role permissions by its ID
	filter := bson.M{"identifier": org_identifier}
	// Find the user document in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Role Permissions"}
		}
		return nil, err
	}

	return &org.RolePermissions, nil
}
