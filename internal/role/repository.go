package role

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
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
}

type repository struct {
	mongodb     *mongo.Database
	writeClient rts.WriteServiceClient
}

func NewRepository(mongodb *mongo.Database, writeClient rts.WriteServiceClient) Repository {

	return repository{mongodb: mongodb, writeClient: writeClient}
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
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
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
	_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
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
	_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
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
	_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// add users
	if len(update_role.AddedUsers) > 0 {

		filter := bson.M{"_id": orgId, "roles._id": roleId}
		update := bson.M{"$push": bson.M{"roles.$.users": bson.M{
			"$each": update_role.AddedUsers,
		}}}
		_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

		for _, userId := range update_role.AddedUsers {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$addToSet": bson.M{"users.$.roles": roleId}}
			_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	// remove users
	if len(update_role.RemovedUser) > 0 {

		filter := bson.M{"_id": orgId, "roles._id": roleId}
		update := bson.M{"$pull": bson.M{"roles.$.users": bson.M{"$in": update_role.RemovedUser}}}
		_, err := r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
		if err != nil {
			return err
		}

		for _, userId := range update_role.RemovedUser {
			filter := bson.M{"_id": orgId, "users._id": userId}
			update := bson.M{"$pull": bson.M{"users.$.roles": roleId}}
			_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update)
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
	result, err := r.mongodb.Collection("organizations").UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// Check if the update operation modified any documents
	if result.ModifiedCount == 0 {
		return err
	}

	// Delete role permission
	filter = bson.M{"_id": orgId}
	update = bson.M{"$pull": bson.M{"role_permissions": bson.M{"role_id": roleId}}}
	// Find the role document in the "organizations" collection
	result, err = r.mongodb.Collection("organizations").UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
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
	// Find the role document in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)
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
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)

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
	count, err := r.mongodb.Collection("organizations").CountDocuments(context.Background(), filter)

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
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)

	// Check if the user was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

func qualifiedTuple(org string, tuple entity.Tuple) entity.Tuple {

	tuple.Object = org + "/" + tuple.Object
	tuple.SubjectId = org + "/" + tuple.SubjectId
	return tuple
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
	err = r.mongodb.Collection("organizations").FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&org)
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
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)

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

	// add permissions
	if len(permission.AddedPermission) > 0 {

		filter := bson.M{"_id": orgId, "role_permissions.role_id": roleId}
		update := bson.M{"$push": bson.M{"role_permissions.$.permissions": bson.M{
			"$each": permission.AddedPermission,
		}}}
		_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	return nil
}
