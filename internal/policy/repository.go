package policy

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
	Get(ctx context.Context, org_id string, id string) (*mongo_entity.Policy, error)
	Create(ctx context.Context, org_id string, policy mongo_entity.Policy) error
	Query(ctx context.Context, org_id string) (*[]mongo_entity.Policy, error)
	Update(ctx context.Context, org_id string, id string, update_user UpdatePolicy) error
	Patch(ctx context.Context, org_id string, id string, patch_user PatchPolicy) error
	Delete(ctx context.Context, org_id string, id string) error
	CheckPolicyExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckPolicyExistsByIdentifier(ctx context.Context, org_id string, identifier string) (bool, error)
	CheckPolicyContentExistsByVersion(ctx context.Context, org_id string, version string) (bool, error)
}

type repository struct {
	mongoClient *mongo.Client
	mongoColl   *mongo.Collection
}

func NewRepository(mongodb *db.MongoDB) Repository {

	orgCollection := mongodb.MongoClient.Database(mongodb.MongoConfig.DBName).Collection(mongodb.MongoConfig.OrganizationCollectionName)

	return repository{mongoClient: mongodb.MongoClient, mongoColl: orgCollection}
}

// Get policy by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (*mongo_entity.Policy, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	policyId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the policy by its ID
	filter := bson.M{"_id": orgId, "policies._id": policyId}
	projection := bson.M{"policies.$": 1}
	// Find the user document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Policy"}
		}
		return nil, err
	}

	return &org.Polices[0], nil
}

// Create new policy.
func (r repository) Create(ctx context.Context, org_id string, policy mongo_entity.Policy) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}
	// Update the APIResources array for the given organization
	filter := bson.M{"_id": orgId}
	update := bson.M{"$push": bson.M{"policies": policy}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil

}

func (r repository) Update(ctx context.Context, org_id string, id string, update_policy UpdatePolicy) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	policyId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": orgId, "policies._id": policyId}
	update := bson.M{"$set": bson.M{}}

	if update_policy.DisplayName != nil && *update_policy.DisplayName != "" {
		update["$set"].(bson.M)["policies.$.display_name"] = *update_policy.DisplayName
	}
	if update_policy.ActiveVersion != nil && *update_policy.ActiveVersion != "" {
		update["$set"].(bson.M)["policies.$.active_version"] = *update_policy.ActiveVersion
	}

	if update_policy.PolicyContent != nil && update_policy.PolicyContent.Version != nil && *update_policy.PolicyContent.Version != "" {
		// Update the specific policy content
		contentFilter := bson.M{
			"_id":                              orgId,
			"policies._id":                     policyId,
			"policies.policy_contents.version": *update_policy.PolicyContent.Version,
		}
		update["$set"].(bson.M)["policies.$.policy_contents.$[elem].policy"] = *update_policy.PolicyContent.Policy

		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{bson.M{"elem.version": *update_policy.PolicyContent.Version}},
		}
		opts := options.Update().SetArrayFilters(arrayFilters).SetUpsert(true)

		_, err = r.mongoColl.UpdateOne(ctx, contentFilter, update, opts)
		if err != nil {
			return err
		}
	} else {
		// If no policy content update is needed, just update the policy
		_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r repository) Patch(ctx context.Context, org_id string, id string, patch_user PatchPolicy) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	policyId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// add roles
	if len(patch_user.AddedPolicies) > 0 {

		filter := bson.M{"_id": orgId, "policies._id": policyId}
		update := bson.M{"$push": bson.M{"policies.$.policy_contents": bson.M{
			"$each": patch_user.AddedPolicies,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	// remove roles
	if len(patch_user.RemovedPolicies) > 0 {

		filter := bson.M{"_id": orgId, "policies._id": policyId}
		update := bson.M{
			"$pull": bson.M{
				"policies.$.policy_contents": bson.M{
					"version": bson.M{"$in": patch_user.RemovedPolicies},
				},
			},
		}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete existing policy.
func (r repository) Delete(ctx context.Context, org_id string, id string) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define filter to find the policy by its ID
	filter := bson.M{"_id": orgId}
	update := bson.M{"$pull": bson.M{"policies": bson.M{"_id": userId}}}
	// Find the policy document in the "organizations" collection
	result, err := r.mongoColl.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// Check if the update operation modified any documents
	if result.ModifiedCount == 0 {
		return err
	}

	// filter = bson.M{"_id": orgId}
	// update = bson.M{"$pull": bson.M{"groups.$[].users": userId}}
	// _, err = r.mongoColl.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	return err
	// }

	// filter = bson.M{"_id": orgId}
	// update = bson.M{"$pull": bson.M{"roles.$[].users": userId}}
	// _, err = r.mongoColl.UpdateOne(ctx, filter, update)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Get all policies.
func (r repository) Query(ctx context.Context, org_id string) (*[]mongo_entity.Policy, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the policy by its ID
	filter := bson.M{"_id": orgId}
	projection := bson.M{"policies.policy_contents": 0}
	// Find the policy document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Policy"}
		}
		return nil, err
	}

	return &org.Polices, nil
}

// Check if policy exists by id.
func (r repository) CheckPolicyExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	policyId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "policies._id": policyId}

	// Search for the policy in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter)

	// Check if the policy was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Check if policy exists by key.
func (r repository) CheckPolicyExistsByIdentifier(ctx context.Context, org_id string, identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": orgId, "policies.identifier": identifier}

	// Search for the policy in the "organizations" collection
	count, err := r.mongoColl.CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Check if policy content exists by version.
func (r repository) CheckPolicyContentExistsByVersion(ctx context.Context, org_id string, version string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": orgId, "policies.policy_contents.version": version}

	// Search for the policy in the "organizations" collection
	count, err := r.mongoColl.CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// // Check if role exists by id.
// func (r repository) CheckRoleExistById(ctx context.Context, org_id string, id string) (bool, error) {

// 	orgId, err := primitive.ObjectIDFromHex(org_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	roleId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return false, err
// 	}

// 	filter := bson.M{"_id": orgId, "roles._id": roleId}

// 	// Search for the role in the "organizations" collection
// 	result := r.mongoColl.FindOne(context.Background(), filter)

// 	// Check if the role was found
// 	if result.Err() == nil {
// 		return true, nil
// 	} else if result.Err() == mongo.ErrNoDocuments {
// 		return false, nil
// 	} else {
// 		return false, result.Err()
// 	}
// }

// // Check if role already assign to user by id.
// func (r repository) CheckRoleAlreadyAssignToUserById(ctx context.Context, org_id string, user_id string, role_id string) (bool, error) {

// 	orgId, err := primitive.ObjectIDFromHex(org_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	userId, err := primitive.ObjectIDFromHex(user_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	roleId, err := primitive.ObjectIDFromHex(role_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	filter := bson.M{"_id": orgId, "users._id": userId}
// 	projection := bson.M{"users.$": 1}
// 	org := mongo_entity.Organization{}
// 	// Search for the role in the "organizations" collection
// 	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
// 	if err != nil {
// 		return false, err
// 	}
// 	user := org.Users[0]
// 	// Check if the role ID exists in the user's Roles field
// 	for _, r := range user.Roles {
// 		if r == roleId {
// 			return true, nil
// 		}
// 	}

// 	// Role ID not found in the user's Roles field
// 	return false, nil
// }

// // Check if group exists by id.
// func (r repository) CheckGroupExistById(ctx context.Context, org_id string, id string) (bool, error) {

// 	orgId, err := primitive.ObjectIDFromHex(org_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	groupId, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return false, err
// 	}

// 	filter := bson.M{"_id": orgId, "groups._id": groupId}

// 	// Search for the group in the "organizations" collection
// 	result := r.mongoColl.FindOne(context.Background(), filter)

// 	// Check if the group was found
// 	if result.Err() == nil {
// 		return true, nil
// 	} else if result.Err() == mongo.ErrNoDocuments {
// 		return false, nil
// 	} else {
// 		return false, result.Err()
// 	}
// }

// // Check if group already assign to user by id.
// func (r repository) CheckGroupAlreadyAssignToUserById(ctx context.Context, org_id string, user_id string, group_id string) (bool, error) {

// 	orgId, err := primitive.ObjectIDFromHex(org_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	userId, err := primitive.ObjectIDFromHex(user_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	groupId, err := primitive.ObjectIDFromHex(group_id)
// 	if err != nil {
// 		return false, err
// 	}

// 	filter := bson.M{"_id": orgId, "users._id": userId}
// 	projection := bson.M{"users.$": 1}
// 	org := mongo_entity.Organization{}
// 	// Search for the role in the "organizations" collection
// 	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
// 	if err != nil {
// 		return false, err
// 	}
// 	user := org.Users[0]
// 	// Check if the group ID exists in the user's Groups field
// 	for _, r := range user.Groups {
// 		if r == groupId {
// 			return true, nil
// 		}
// 	}

// 	// Group ID not found in the user's Groups field
// 	return false, nil
// }
