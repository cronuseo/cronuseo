package group

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
	Get(ctx context.Context, org_id string, id string) (*GroupResponse, error)
	Query(ctx context.Context, org_id string) (*[]mongo_entity.Group, error)
	Create(ctx context.Context, org_id string, group mongo_entity.Group) error
	Update(ctx context.Context, org_id string, id string, update_group UpdateGroup) error
	Patch(ctx context.Context, org_id string, id string, patch_group PatchGroup) error
	Delete(ctx context.Context, org_id string, id string) error
	CheckGroupExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckGroupExistsByIdentifier(ctx context.Context, org_id string, key string) (bool, error)
	CheckRoleExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckRoleAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, role_id string) (bool, error)
	CheckUserExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckUserAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, user_id string) (bool, error)
	CheckPolicyExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckPolicyAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, policy_id string) (bool, error)
}

type repository struct {
	mongoClient *mongo.Client
	mongoColl   *mongo.Collection
}

func NewRepository(mongodb *db.MongoDB) Repository {

	orgCollection := mongodb.MongoClient.Database(mongodb.MongoConfig.DBName).Collection(mongodb.MongoConfig.OrganizationCollectionName)

	return repository{mongoClient: mongodb.MongoClient, mongoColl: orgCollection}
}

// Get group by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (*GroupResponse, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the group by its ID
	filter := bson.M{"_id": orgId, "groups._id": groupId}
	projection := bson.M{"groups.$": 1}
	// Find the group document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Group"}
		}
		return nil, err
	}

	group := org.Groups[0]
	assignedUsers, err := r.resolveAssignedUsers(ctx, orgId, group.Users)
	assignedRoles, err := r.resolveAssignedRoles(ctx, orgId, group.Roles)
	assignedPolicies, err := r.resolveAssignedPolicies(ctx, orgId, group.Policies)
	if err != nil {
		return nil, err
	}
	roleResponse := GroupResponse{
		ID:          group.ID,
		Identifier:  group.Identifier,
		DisplayName: group.DisplayName,
		Users:       assignedUsers,
		Roles:       assignedRoles,
		Policies:    assignedPolicies,
	}
	return &roleResponse, nil
}

// Create new group.
func (r repository) Create(ctx context.Context, org_id string, group mongo_entity.Group) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}
	// Update the APIResources array for the given organization
	filter := bson.M{"_id": orgId}
	update := bson.M{"$push": bson.M{"groups": group}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// add roles
	if len(group.Roles) > 0 {

		for _, roleId := range group.Roles {
			filter := bson.M{"_id": orgId, "roles._id": roleId}
			update := bson.M{"$addToSet": bson.M{"roles.$.groups": group.ID}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// add users
	if len(group.Users) > 0 {

		for _, userId := range group.Users {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$addToSet": bson.M{"users.$.groups": group.ID}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}

	}

	return nil

}

func (r repository) Update(ctx context.Context, org_id string, id string, update_group UpdateGroup) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": orgId, "groups._id": groupId}
	update := bson.M{"$set": bson.M{}}
	if update_group.DisplayName != nil && *update_group.DisplayName != "" {
		update["$set"].(bson.M)["groups.$.first_name"] = *update_group.DisplayName
	}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Patch(ctx context.Context, org_id string, id string, patch_group PatchGroup) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// add roles
	if len(patch_group.AddedRoles) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$push": bson.M{"groups.$.roles": bson.M{
			"$each": patch_group.AddedRoles,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, roleId := range patch_group.AddedRoles {
			filter := bson.M{"_id": orgId, "roles._id": roleId}
			update := bson.M{"$addToSet": bson.M{"roles.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}

	}

	// remove roles
	if len(patch_group.RemovedRoles) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$pull": bson.M{"groups.$.roles": bson.M{"$in": patch_group.RemovedRoles}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, roleId := range patch_group.RemovedRoles {
			filter := bson.M{"_id": orgId, "roles._id": roleId}
			update := bson.M{"$pull": bson.M{"roles.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// add users
	if len(patch_group.AddedUsers) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$push": bson.M{"groups.$.users": bson.M{
			"$each": patch_group.AddedUsers,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, userId := range patch_group.AddedUsers {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$addToSet": bson.M{"users.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}

	}

	// remove users
	if len(patch_group.RemovedUsers) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$pull": bson.M{"groups.$.users": bson.M{"$in": patch_group.RemovedUsers}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, userId := range patch_group.RemovedUsers {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$pull": bson.M{"users.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// add policies
	if len(patch_group.AddedPolicies) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$push": bson.M{"groups.$.policies": bson.M{
			"$each": patch_group.AddedPolicies,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	// remove policies
	if len(patch_group.RemovedPolicies) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$pull": bson.M{"groups.$.policies": bson.M{"$in": patch_group.RemovedPolicies}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete existing group.
func (r repository) Delete(ctx context.Context, org_id string, id string) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define filter to find the group by its ID
	filter := bson.M{"_id": orgId}
	update := bson.M{"$pull": bson.M{"groups": bson.M{"_id": groupId}}}
	// Find the group document in the "organizations" collection
	result, err := r.mongoColl.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// Check if the update operation modified any documents
	if result.ModifiedCount == 0 {
		return err
	}

	filter = bson.M{"_id": orgId}
	update = bson.M{"$pull": bson.M{"roles.$[].groups": groupId}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	filter = bson.M{"_id": orgId}
	update = bson.M{"$pull": bson.M{"users.$[].groups": groupId}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// Get all groups.
func (r repository) Query(ctx context.Context, org_id string) (*[]mongo_entity.Group, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the group by its ID
	filter := bson.M{"_id": orgId}
	projection := bson.M{"groups.roles": 0, "groups.users": 0}
	// Find the group document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Group"}
		}
		return nil, err
	}

	return &org.Groups, nil
}

// Check if group exists by id.
func (r repository) CheckGroupExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "groups._id": groupId}

	// Search for the group in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter)

	// Check if the group was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Check if group exists by key.
func (r repository) CheckGroupExistsByIdentifier(ctx context.Context, org_id string, identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": orgId, "groups.identifier": identifier}

	// Search for the group in the "organizations" collection
	count, err := r.mongoColl.CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Check if role exists by id.
func (r repository) CheckRoleExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	roleId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "roles._id": roleId}

	// Search for the role in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter)

	// Check if the role was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Check if role already assign to group by id.
func (r repository) CheckRoleAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, role_id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	groupId, err := primitive.ObjectIDFromHex(group_id)
	if err != nil {
		return false, err
	}

	roleId, err := primitive.ObjectIDFromHex(role_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "groups._id": groupId}
	projection := bson.M{"groups.$": 1}
	org := mongo_entity.Organization{}
	// Search for the role in the "organizations" collection
	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
	if err != nil {
		return false, err
	}
	group := org.Groups[0]
	// Check if the role ID exists in the group's Roles field
	for _, r := range group.Roles {
		if r == roleId {
			return true, nil
		}
	}

	// Role ID not found in the group's Roles field
	return false, nil
}

// Check if user exists by id.
func (r repository) CheckUserExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "users._id": userId}

	// Search for the user in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter)

	// Check if the user was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Check if role already assign to group by id.
func (r repository) CheckUserAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, user_id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	groupId, err := primitive.ObjectIDFromHex(group_id)
	if err != nil {
		return false, err
	}

	userId, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "groups._id": groupId}
	projection := bson.M{"groups.$": 1}
	org := mongo_entity.Organization{}
	// Search for the user in the "organizations" collection
	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
	if err != nil {
		return false, err
	}
	group := org.Groups[0]
	// Check if the user ID exists in the group's Users field
	for _, r := range group.Users {
		if r == userId {
			return true, nil
		}
	}

	// User ID not found in the group's Users field
	return false, nil
}

// Check if policy exists by id.
func (r repository) CheckPolicyExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	groupId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "policies._id": groupId}

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

// Check if policy already assign to group by id.
func (r repository) CheckPolicyAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, policy_id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	userId, err := primitive.ObjectIDFromHex(group_id)
	if err != nil {
		return false, err
	}

	policyId, err := primitive.ObjectIDFromHex(policy_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "groups._id": userId}
	projection := bson.M{"groups.$": 1}
	org := mongo_entity.Organization{}
	// Search for the role in the "organizations" collection
	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
	if err != nil {
		return false, err
	}
	group := org.Groups[0]
	// Check if the policy ID exists in the user's policies field
	for _, r := range group.Policies {
		if r == policyId {
			return true, nil
		}
	}

	// policy ID not found in the user's policies field
	return false, nil
}

func (r repository) resolveAssignedUsers(ctx context.Context, orgId primitive.ObjectID, userIDs []primitive.ObjectID) ([]mongo_entity.AssignedUser, error) {

	filter := bson.M{"_id": orgId, "users._id": bson.M{"$in": userIDs}}

	cursor, err := r.mongoColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []mongo_entity.AssignedUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r repository) resolveAssignedRoles(ctx context.Context, orgId primitive.ObjectID, roleIDs []primitive.ObjectID) ([]mongo_entity.AssignedRole, error) {

	filter := bson.M{"_id": orgId, "roles._id": bson.M{"$in": roleIDs}}

	cursor, err := r.mongoColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []mongo_entity.AssignedRole
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r repository) resolveAssignedPolicies(ctx context.Context, orgId primitive.ObjectID, policyIDs []primitive.ObjectID) ([]mongo_entity.AssignedPolicy, error) {

	filter := bson.M{"_id": orgId, "groups._id": bson.M{"$in": policyIDs}}

	cursor, err := r.mongoColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var policies []mongo_entity.AssignedPolicy
	if err := cursor.All(ctx, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}
