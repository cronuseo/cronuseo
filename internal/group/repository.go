package group

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
	Get(ctx context.Context, org_id string, id string) (*mongo_entity.Group, error)
	Query(ctx context.Context, org_id string) (*[]mongo_entity.Group, error)
	Create(ctx context.Context, org_id string, group mongo_entity.Group) error
	Update(ctx context.Context, org_id string, id string, update_group UpdateGroup) error
	// Patch(ctx context.Context, org_id string, id string, req GroupPatchRequest) error
	Delete(ctx context.Context, org_id string, id string) error
	CheckGroupExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckGroupExistsByIdentifier(ctx context.Context, org_id string, key string) (bool, error)
	CheckRoleExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckRoleAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, role_id string) (bool, error)
	CheckUserExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckUserAlreadyAssignToGroupById(ctx context.Context, org_id string, group_id string, user_id string) (bool, error)
}

type repository struct {
	mongoClient   *mongo.Client
	mongoDBConfig util.MongoDBConfig
	mongoColl     *mongo.Collection
}

func NewRepository(mongoClient *mongo.Client, mongoDBConfig util.MongoDBConfig) Repository {

	orgCollection := mongoClient.Database(mongoDBConfig.DBName).Collection(mongoDBConfig.OrganizationCollectionName)

	return repository{mongoClient: mongoClient, mongoDBConfig: mongoDBConfig, mongoColl: orgCollection}
}

// Get group by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (*mongo_entity.Group, error) {

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

	return &org.Groups[0], nil
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

	// add roles
	if len(update_group.AddedRoles) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$push": bson.M{"groups.$.roles": bson.M{
			"$each": update_group.AddedRoles,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, roleId := range update_group.AddedRoles {
			filter := bson.M{"_id": orgId, "roles._id": roleId}
			update := bson.M{"$addToSet": bson.M{"roles.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}

	}

	// remove roles
	if len(update_group.RemovedRoles) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$pull": bson.M{"groups.$.roles": bson.M{"$in": update_group.RemovedRoles}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, roleId := range update_group.RemovedRoles {
			filter := bson.M{"_id": orgId, "roles._id": roleId}
			update := bson.M{"$pull": bson.M{"roles.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// add users
	if len(update_group.AddedUsers) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$push": bson.M{"groups.$.users": bson.M{
			"$each": update_group.AddedUsers,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, userId := range update_group.AddedUsers {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$addToSet": bson.M{"users.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}

	}

	// remove roles
	if len(update_group.RemovedUsers) > 0 {

		filter := bson.M{"_id": orgId, "groups._id": groupId}
		update := bson.M{"$pull": bson.M{"groups.$.users": bson.M{"$in": update_group.RemovedUsers}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, userId := range update_group.RemovedUsers {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$pull": bson.M{"users.$.groups": groupId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
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
