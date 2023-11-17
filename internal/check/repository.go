package check

import (
	"context"

	db "github.com/shashimalcse/cronuseo/internal/db/mongo"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	ValidateAPIKey(ctx context.Context, org_identifier string, apiKey string) (bool, error)
	GetUserRoles(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error)
	GetRolePermissions(ctx context.Context, org_identifier string, role_ids []primitive.ObjectID) (*[]mongo_entity.Permission, error)
	GetGroupRoles(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error)
	GetUserGroups(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error)
	GetUserPolicies(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error)
	GetGroupPolicies(ctx context.Context, org_identifier string, group_ids []primitive.ObjectID) (*[]primitive.ObjectID, error)
	GetUserProperties(ctx context.Context, org_identifier string, identifier string) (*map[string]interface{}, error)
	GetActivePolicyVersionContents(ctx context.Context, org_identifier string, policy_ids []primitive.ObjectID) (map[string]string, error)
}

type repository struct {
	mongoClient *mongo.Client
	mongoColl   *mongo.Collection
}

func NewRepository(mongodb *db.MongoDB) Repository {

	orgCollection := mongodb.MongoClient.Database(mongodb.MongoConfig.DBName).Collection(mongodb.MongoConfig.OrganizationCollectionName)

	return repository{mongoClient: mongodb.MongoClient, mongoColl: orgCollection}
}

func (r repository) ValidateAPIKey(ctx context.Context, org_identifier string, apiKey string) (bool, error) {

	filter := bson.M{"identifier": org_identifier, "api_key": apiKey}

	// Search for the resource in the "organizations" collection
	count, err := r.mongoColl.CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r repository) GetUserRoles(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error) {

	// Define filter to find the user by its ID
	filter := bson.M{"identifier": org_identifier, "users.identifier": identifier}
	projection := bson.M{"users.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))

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

func (r repository) GetGroupRoles(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error) {

	// Define filter to find the user by its ID
	filter := bson.M{"identifier": org_identifier, "users.identifier": identifier}
	projection := bson.M{"users.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))

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
	err := r.mongoColl.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&org2)
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

func (r repository) GetUserGroups(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error) {

	// Define filter to find the user by its ID
	filter := bson.M{"identifier": org_identifier, "users.identifier": identifier}
	projection := bson.M{"users.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "User"}
		}
		return nil, err
	}

	return &org.Users[0].Groups, nil
}

func (r repository) GetUserPolicies(ctx context.Context, org_identifier string, identifier string) (*[]primitive.ObjectID, error) {

	// Define filter to find the user by its ID
	filter := bson.M{"identifier": org_identifier, "users.identifier": identifier}
	projection := bson.M{"users.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "User"}
		}
		return nil, err
	}

	return &org.Users[0].Policies, nil
}

func (r repository) GetGroupPolicies(ctx context.Context, org_identifier string, group_ids []primitive.ObjectID) (*[]primitive.ObjectID, error) {

	// Define filter to find the user by its ID
	filter := bson.M{"identifier": org_identifier, "groups._id": bson.M{"$in": group_ids}}
	projection := bson.M{"groups.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "User"}
		}
		return nil, err
	}

	return &org.Groups[0].Policies, nil
}

func (r repository) GetRolePermissions(ctx context.Context, org_identifier string, role_ids []primitive.ObjectID) (*[]mongo_entity.Permission, error) {

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"identifier": org_identifier}}},
		{{Key: "$project", Value: bson.M{
			"roles": bson.M{
				"$filter": bson.M{
					"input": "$roles",
					"as":    "role",
					"cond":  bson.M{"$in": []interface{}{"$$role._id", role_ids}},
				},
			},
		}}},
	}

	// Aggregate the query results
	cursor, err := r.mongoColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Assuming only one organization will match, decode the first result
	var org struct {
		Roles []mongo_entity.Role `bson:"roles"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&org); err != nil {
			return nil, err
		}
	} else {
		return nil, &util.NotFoundError{Path: "Organization not found"}
	}

	// Initialize a slice to store permissions
	var permissions []mongo_entity.Permission
	for _, role := range org.Roles {
		permissions = append(permissions, role.Permissions...)
	}

	// Return the collected permissions
	return &permissions, nil
}

func (r repository) GetActivePolicyVersionContents(ctx context.Context, org_identifier string, policy_ids []primitive.ObjectID) (map[string]string, error) {

	// Filter to find documents with the specified policy IDs
	filter := bson.M{
		"identifier": org_identifier,
		"policies": bson.M{
			"$elemMatch": bson.M{
				"_id": bson.M{"$in": policy_ids},
			},
		},
	}

	// Projection to extract only the relevant policies
	projection := bson.M{"policies": 1}

	cursor, err := r.mongoColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	activePolicies := make(map[string]string)
	for cursor.Next(ctx) {
		var doc struct {
			Policies []struct {
				ID             primitive.ObjectID `bson:"_id"`
				ActiveVersion  string             `bson:"active_version"`
				PolicyContents []struct {
					Version string `bson:"version"`
					Policy  string `bson:"policy"`
				} `bson:"policy_contents"`
			} `bson:"policies"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		for _, policy := range doc.Policies {
			if contains(policy_ids, policy.ID) {
				for _, content := range policy.PolicyContents {
					if content.Version == policy.ActiveVersion {
						activePolicies[policy.ID.Hex()] = content.Policy
						break
					}
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return activePolicies, nil
}

func (r repository) GetUserProperties(ctx context.Context, org_identifier string, identifier string) (*map[string]interface{}, error) {

	filter := bson.M{"identifier": org_identifier, "users.identifier": identifier}
	projection := bson.M{"users.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "User"}
		}
		return nil, err
	}

	return &org.Users[0].UserProperties, nil
}

func contains(slice []primitive.ObjectID, item primitive.ObjectID) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
