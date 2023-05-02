package role

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
	Get(ctx context.Context, org_id string, id string) (*mongo_entity.Role, error)
	Query(ctx context.Context, org_id string) (*[]mongo_entity.Role, error)
	// QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]entity.Role, error)
	Create(ctx context.Context, org_id string, user mongo_entity.Role) error
	Update(ctx context.Context, org_id string, id string, update_role UpdateRole) error
	Delete(ctx context.Context, org_id string, id string) error
	CheckRoleExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckRoleExistsByIdentifier(ctx context.Context, org_id string, key string) (bool, error)
	CheckUserExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckUserAlreadyAssignToRoleById(ctx context.Context, org_id string, role_id string, user_id string) (bool, error)
	CheckResourceActionExists(ctx context.Context, org_id string, resource_identifier string, action_identifier string) (bool, error)
	PatchPermission(ctx context.Context, org_id string, role_id string, permission PatchRolePermission) error
	CheckPermissionExists(ctx context.Context, org_id string, role_id string, resource_identifier string, action_identifier string) (bool, error)
	GetPermissions(ctx context.Context, org_id string, role_id string) (*[]mongo_entity.Permission, error)
	CheckGroupExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckGroupAlreadyAssignToRoleById(ctx context.Context, org_id string, role_id string, group_id string) (bool, error)
}

type repository struct {
	mongoClient *mongo.Client
	mongoColl   *mongo.Collection
}

func NewRepository(mongodb *db.MongoDB) Repository {

	orgCollection := mongodb.MongoClient.Database(mongodb.MongoConfig.DBName).Collection(mongodb.MongoConfig.OrganizationCollectionName)

	return repository{mongoClient: mongodb.MongoClient, mongoColl: orgCollection}
}

// Get role by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (*mongo_entity.Role, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	roleId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the role by its ID
	filter := bson.M{"_id": orgId, "roles._id": roleId}
	projection := bson.M{"roles.$": 1}
	// Find the role document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}
	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Role"}
		}
		return nil, err
	}

	return &org.Roles[0], nil
}

// Create new role.
func (r repository) Create(ctx context.Context, org_id string, user mongo_entity.Role) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}
	// Update the APIResources array for the given organization
	filter := bson.M{"_id": orgId}
	update := bson.M{"$push": bson.M{"roles": user}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// Set role permissions
	rolePermission := mongo_entity.RolePermission{
		RoleID:      user.ID,
		Permissions: []mongo_entity.Permission{},
	}
	filter = bson.M{"_id": orgId}
	update = bson.M{"$push": bson.M{"role_permissions": rolePermission}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil

}

// Update role.
func (r repository) Update(ctx context.Context, org_id string, id string, update_role UpdateRole) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	roleId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": orgId, "roles._id": roleId}
	update := bson.M{"$set": bson.M{}}

	if update_role.DisplayName != nil && *update_role.DisplayName != "" {
		update["$set"].(bson.M)["roles.$.display_name"] = *update_role.DisplayName
	}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// add users
	if len(update_role.AddedUsers) > 0 {

		filter := bson.M{"_id": orgId, "roles._id": roleId}
		update := bson.M{"$push": bson.M{"roles.$.users": bson.M{
			"$each": update_role.AddedUsers,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, userId := range update_role.AddedUsers {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$addToSet": bson.M{"users.$.roles": roleId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// remove users
	if len(update_role.RemovedUser) > 0 {

		filter := bson.M{"_id": orgId, "roles._id": roleId}
		update := bson.M{"$pull": bson.M{"roles.$.users": bson.M{"$in": update_role.RemovedUser}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, userId := range update_role.RemovedUser {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$pull": bson.M{"users.$.roles": roleId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// add groups
	if len(update_role.AddedGroups) > 0 {

		filter := bson.M{"_id": orgId, "roles._id": roleId}
		update := bson.M{"$push": bson.M{"roles.$.groups": bson.M{
			"$each": update_role.AddedGroups,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, groupId := range update_role.AddedGroups {
			filter := bson.M{"_id": orgId, "groups._id": groupId}
			update := bson.M{"$addToSet": bson.M{"groups.$.roles": roleId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}

	}

	// remove groups
	if len(update_role.RemovedGroups) > 0 {

		filter := bson.M{"_id": orgId, "roles._id": roleId}
		update := bson.M{"$pull": bson.M{"roles.$.groups": bson.M{"$in": update_role.RemovedGroups}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, groupId := range update_role.RemovedGroups {
			filter := bson.M{"_id": orgId, "groups._id": groupId}
			update := bson.M{"$pull": bson.M{"groups.$.roles": roleId}}
			_, err = r.mongoColl.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete role.
func (r repository) Delete(ctx context.Context, org_id string, id string) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	roleId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define filter to find the role by its ID
	filter := bson.M{"_id": orgId}
	update := bson.M{"$pull": bson.M{"roles": bson.M{"_id": roleId}}}
	// Find the role document in the "organizations" collection
	result, err := r.mongoColl.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// Check if the update operation modified any documents
	if result.ModifiedCount == 0 {
		return err
	}

	filter = bson.M{"_id": orgId}
	update = bson.M{"$pull": bson.M{"groups.$[].roles": roleId}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	filter = bson.M{"_id": orgId}
	update = bson.M{"$pull": bson.M{"users.$[].roles": roleId}}
	_, err = r.mongoColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// Delete role permission
	filter = bson.M{"_id": orgId}
	update = bson.M{"$pull": bson.M{"role_permissions": bson.M{"role_id": roleId}}}
	// Find the role document in the "organizations" collection
	result, err = r.mongoColl.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// Check if the update operation modified any documents
	if result.ModifiedCount == 0 {
		return err
	}

	return nil
}

// Query roles.
func (r repository) Query(ctx context.Context, org_id string) (*[]mongo_entity.Role, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the role by its ID
	filter := bson.M{"_id": orgId}
	projection := bson.M{"roles.groups": 0, "roles.users": 0}
	// Find the role document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Role"}
		}
		return nil, err
	}

	return &org.Roles, nil
}

// Query roles by user id.
// func (r repository) QueryByUserID(ctx context.Context, org_id string, user_id string, filter Filter) ([]entity.Role, error) {
// 	roles := []entity.Role{}
// 	name := filter.Name + "%"
// 	err := r.db.Select(&roles, "SELECT org_role.id, org_role.role_id, org_role.role_key, org_role.name, org_role.org_id, org_role.created_at, org_role.updated_at FROM org_role INNER JOIN user_role ON org_role.role_id = user_role.role_id WHERE org_role.org_id = $1 AND user_role.user_id = $2 AND org_role.name LIKE $3 AND org_role.id > $4 ORDER BY org_role.id LIMIT $5", org_id, user_id, name, filter.Cursor, filter.Limit)
// 	return roles, err
// }

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

// Check if role exists by key.
func (r repository) CheckRoleExistsByIdentifier(ctx context.Context, org_id string, identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": orgId, "roles.identifier": identifier}

	// Search for the user in the "organizations" collection
	count, err := r.mongoColl.CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// Check if user exists by id.
func (r repository) CheckUserExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	roleId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "users._id": roleId}

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

// check user already added to role
func (r repository) CheckUserAlreadyAssignToRoleById(ctx context.Context, org_id string, role_id string, user_id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	userId, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		return false, err
	}

	roleId, err := primitive.ObjectIDFromHex(role_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "roles._id": roleId}
	projection := bson.M{"roles.$": 1}
	org := mongo_entity.Organization{}
	// Search for the role in the "organizations" collection
	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
	if err != nil {
		return false, err
	}
	role := org.Roles[0]
	// Check if the user ID exists in the role's Users field
	for _, r := range role.Users {
		if r == userId {
			return true, nil
		}
	}

	// User ID not found in the role's Roles field
	return false, nil
}

// check user already added to role
func (r repository) CheckResourceActionExists(ctx context.Context, org_id string, resource_identifier string, action_identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "resources.identifier": resource_identifier, "resources.actions.identifier": action_identifier}
	result := r.mongoColl.FindOne(context.Background(), filter)

	// Check if the resource was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Patch permissions.
func (r repository) PatchPermission(ctx context.Context, org_id string, role_id string, permission PatchRolePermission) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	roleId, err := primitive.ObjectIDFromHex(role_id)
	if err != nil {
		return err
	}

	// remove permissions
	if len(permission.RemovedPermission) > 0 {

		filter := bson.M{"_id": orgId, "role_permissions.role_id": roleId}
		update := bson.M{"$pull": bson.M{"role_permissions.$.permissions": bson.M{"$in": permission.RemovedPermission}}}
		_, err := r.mongoColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}
	}

	// add permissions
	if len(permission.AddedPermission) > 0 {

		filter := bson.M{"_id": orgId, "role_permissions.role_id": roleId}
		update := bson.M{"$push": bson.M{"role_permissions.$.permissions": bson.M{
			"$each": permission.AddedPermission,
		}}}
		_, err = r.mongoColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}

// check user already added to role
func (r repository) CheckPermissionExists(ctx context.Context, org_id string, role_id string, resource_identifier string, action_identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	roleId, err := primitive.ObjectIDFromHex(role_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "role_permissions.role_id": roleId, "role_permissions.permissions.resource": resource_identifier, "role_permissions.permissions.action": action_identifier}
	result := r.mongoColl.FindOne(context.Background(), filter)

	// Check if the resource was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Get all resources.
func (r repository) GetPermissions(ctx context.Context, org_id string, role_id string) (*[]mongo_entity.Permission, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	roleId, err := primitive.ObjectIDFromHex(role_id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the permission by role ID
	filter := bson.M{"_id": orgId, "role_permissions.role_id": roleId}
	projection := bson.M{"role_permissions.$": 1}
	// Find the permission document in the "organizations" collection
	result := r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Role Permission"}
		}
		return nil, err
	}

	return &org.RolePermissions[0].Permissions, nil
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

// Check if group already assign to user by id.
func (r repository) CheckGroupAlreadyAssignToRoleById(ctx context.Context, org_id string, role_id string, group_id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	roleId, err := primitive.ObjectIDFromHex(role_id)
	if err != nil {
		return false, err
	}

	groupId, err := primitive.ObjectIDFromHex(group_id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "roles._id": roleId}
	projection := bson.M{"roles.$": 1}
	org := mongo_entity.Organization{}
	// Search for the role in the "organizations" collection
	err = r.mongoColl.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
	if err != nil {
		return false, err
	}
	user := org.Users[0]
	// Check if the group ID exists in the user's Groups field
	for _, r := range user.Groups {
		if r == groupId {
			return true, nil
		}
	}

	// Group ID not found in the user's Groups field
	return false, nil
}
