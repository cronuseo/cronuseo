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
	GetRolePermissions(ctx context.Context, org_identifier string, role_ids []primitive.ObjectID) (*[]mongo_entity.Permission, error)
	GetCheckDetails(ctx context.Context, org_identifier string, identifier string) (CheckDetails, error)
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

func (r repository) GetCheckDetails(ctx context.Context, org_identifier string, identifier string) (CheckDetails, error) {

	filter := bson.M{"identifier": org_identifier, "users.identifier": identifier}
	projection := bson.M{"users.$": 1, "groups": 1}

	// Find the user and groups in the "organizations" collection
	var org mongo_entity.Organization
	err := r.mongoColl.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return CheckDetails{}, &util.NotFoundError{Path: "User"}
		}
		return CheckDetails{}, err
	}

	if len(org.Users) == 0 {
		return CheckDetails{}, nil
	}

	// Create a map to store the unique role IDs
	roleIDMap := make(map[primitive.ObjectID]struct{})
	policyIDMap := make(map[primitive.ObjectID]struct{})
	groupIDs := make(map[primitive.ObjectID]struct{})

	for _, groupID := range org.Users[0].Groups {
		groupIDs[groupID] = struct{}{}
	}
	for _, policyID := range org.Users[0].Policies {
		policyIDMap[policyID] = struct{}{}
	}

	for _, group := range org.Groups {
		if _, exists := groupIDs[group.ID]; exists {
			for _, roleID := range group.Roles {
				roleIDMap[roleID] = struct{}{}
			}
			for _, policyID := range group.Policies {
				policyIDMap[policyID] = struct{}{}
			}
		}
	}

	var roleIDs []primitive.ObjectID
	roleIDs = append(roleIDs, org.Users[0].Roles...)
	for roleID := range roleIDMap {
		roleIDs = append(roleIDs, roleID)
	}

	var policyIDs []primitive.ObjectID
	for policyID := range policyIDMap {
		policyIDs = append(policyIDs, policyID)
	}

	return CheckDetails{
		Roles:          roleIDs,
		Policies:       policyIDs,
		UserProperties: org.Users[0].UserProperties,
	}, nil
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

func contains(slice []primitive.ObjectID, item primitive.ObjectID) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
