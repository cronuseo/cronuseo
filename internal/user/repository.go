package user

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
	Get(ctx context.Context, org_id string, id string) (*mongo_entity.User, error)
	// Query(ctx context.Context, org_id string, filter Filter) (*[]mongo_entity.User, error)
	Create(ctx context.Context, org_id string, user mongo_entity.User) error
	// Update(ctx context.Context, org_id string, user entity.User) error
	// Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) error
	// Delete(ctx context.Context, org_id string, id string) error
	CheckUserExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckUserExistsByIdentifier(ctx context.Context, org_id string, key string) (bool, error)
}

type repository struct {
	mongodb *mongo.Database
}

func NewRepository(mongodb *mongo.Database) Repository {

	return repository{mongodb: mongodb}
}

// Get user by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (*mongo_entity.User, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the resource by its ID
	filter := bson.M{"_id": orgId, "users._id": userId}
	projection := bson.M{"users.$": 1}
	// Find the resource document in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &util.NotFoundError{Path: "Resource"}
		}
		return nil, err
	}

	return &org.Users[0], nil
}

// Create new user.
func (r repository) Create(ctx context.Context, org_id string, user mongo_entity.User) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}
	// Update the APIResources array for the given organization
	filter := bson.M{"_id": orgId}
	update := bson.M{"$push": bson.M{"users": user}}
	_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil

}

// Update user.
// func (r repository) Update(ctx context.Context, org_id string, user entity.User) error {

// 	stmt, err := r.db.Prepare("UPDATE org_user SET firstname = $1, lastname = $2 WHERE org_id = $3 AND user_id = $4")
// 	if err != nil {
// 		return err
// 	}

// 	if _, err = stmt.Exec(user.FirstName, user.LastName, org_id, user.ID); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Delete user.
// func (r repository) Delete(ctx context.Context, org_id string, id string) error {

// 	tx, err := r.db.DB.Begin()

// 	if err != nil {
// 		return err
// 	}

// 	// Delete roles assigned to the user.
// 	{
// 		stmt, err := tx.Prepare("DELETE FROM user_role WHERE user_id = $1")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		if _, err = stmt.Exec(id); err != nil {
// 			return err
// 		}
// 	}

// 	// Delete user.
// 	{
// 		stmt, err := tx.Prepare("DELETE FROM org_user WHERE org_id = $1 AND user_id = $2")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		if _, err = stmt.Exec(org_id, id); err != nil {
// 			return err
// 		}
// 	}

// 	{
// 		if err := tx.Commit(); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // Get all users.
// func (r repository) Query(ctx context.Context, org_id string, filter Filter) ([]entity.User, error) {

// 	users := []entity.User{}
// 	name := filter.Name + "%"
// 	err := r.db.Select(&users, "SELECT * FROM org_user WHERE org_id = $1 AND username LIKE $2 AND id > $3 ORDER BY id LIMIT $4", org_id, name, filter.Cursor, filter.Limit)
// 	return users, err
// }

// // Patch user.
// func (r repository) Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) error {
// 	tx, err := r.db.DB.Begin()

// 	if err != nil {
// 		return err
// 	}

// 	{
// 		for _, operation := range req.Operations {
// 			switch operation.Path {
// 			case "roles":
// 				{
// 					switch operation.Operation {
// 					case "add":
// 						if len(operation.Values) > 0 {
// 							stmt, err := tx.Prepare("INSERT INTO user_role(user_id,role_id) VALUES($1, $2)")
// 							if err != nil {
// 								return err
// 							}
// 							defer stmt.Close()
// 							for _, roleId := range operation.Values {
// 								exists, err := r.isRoleAlreadyAssigned(roleId.Value, id)
// 								if exists {
// 									continue
// 								}
// 								if err != nil {
// 									return err
// 								}
// 								_, err = stmt.Exec(id, roleId.Value)
// 								if err != nil {
// 									return err
// 								}
// 							}
// 						}
// 					case "remove":
// 						if len(operation.Values) > 0 {
// 							stmt, err := tx.Prepare("DELETE FROM user_role WHERE user_id = $1 AND role_id = $2")
// 							if err != nil {
// 								return err
// 							}
// 							defer stmt.Close()
// 							for _, roleId := range operation.Values {
// 								exists, err := r.isRoleAlreadyAssigned(roleId.Value, id)
// 								if !exists {
// 									continue
// 								}
// 								if err != nil {
// 									return err
// 								}
// 								_, err = stmt.Exec(id, roleId.Value)
// 								if err != nil {
// 									return err
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	{
// 		if err := tx.Commit(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // Check if role is already assigned to the user.
// func (r repository) isRoleAlreadyAssigned(roleId string, userId string) (bool, error) {

// 	exists := false
// 	err := r.db.QueryRow("SELECT exists (SELECT role_id FROM user_role WHERE role_id = $1 AND user_id = $2)", roleId, userId).Scan(&exists)
// 	return exists, err
// }

// // Get user by id.
// func (r repository) ExistByID(ctx context.Context, id string) (bool, error) {

// 	exists := false
// 	err := r.db.QueryRow("SELECT exists (SELECT user_id FROM org_user WHERE user_id = $1)", id).Scan(&exists)
// 	return exists, err
// }

// // Get user by username.
// func (r repository) ExistByKey(ctx context.Context, username string) (bool, error) {

// 	exists := false
// 	err := r.db.QueryRow("SELECT exists (SELECT user_id FROM org_user WHERE username = $1)", username).Scan(&exists)
// 	return exists, err
// }

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

// Check if user exists by key.
func (r repository) CheckUserExistsByIdentifier(ctx context.Context, org_id string, identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": orgId, "users.username": identifier}

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
