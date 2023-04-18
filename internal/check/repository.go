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
	GetGroupRoles(ctx context.Context, org_identifier string, username string) (*[]primitive.ObjectID, error)
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

func (r repository) GetGroupRoles(ctx context.Context, org_identifier string, username string) (*[]primitive.ObjectID, error) {

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

	groupIds := org.Users[0].Groups

	if len(groupIds) == 0 {
		return &[]primitive.ObjectID{}, nil
	}

	filter = bson.M{"identifier": org_identifier, "groups._id": bson.M{"$in": groupIds}}

	// create a projection to include only the role IDs
	projection = bson.M{"groups.$": 1, "_id": 0}

	var org2 mongo_entity.Organization
	err := r.mongodb.Collection("organizations").FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&org2)
	if err != nil {
		return nil, err
	}
	var roleIDs []primitive.ObjectID
	// extract the role IDs from the groups that match the specified IDs
	for _, group := range org2.Groups {
		for _, groupID := range groupIds {
			if group.ID == groupID {
				roleIDs = append(roleIDs, group.Roles...)
				break
			}
		}
	}
	return &roleIDs, nil
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
