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
	Update(ctx context.Context, org_id string, id string, update_user UpdateUser) error
	// Patch(ctx context.Context, org_id string, id string, req UserPatchRequest) error
	Delete(ctx context.Context, org_id string, id string) error
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

func (r repository) Update(ctx context.Context, org_id string, id string, update_user UpdateUser) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": orgId, "users._id": userId}
	update := bson.M{"$set": bson.M{"user.$.first_name": update_user.FirstName, "user.$.last_name": update_user.LastName}}
	_, err = r.mongodb.Collection("organizations").UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// Delete existing user.
func (r repository) Delete(ctx context.Context, org_id string, id string) error {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define filter to find the user by its ID
	filter := bson.M{"_id": orgId}
	update := bson.M{"$pull": bson.M{"users": bson.M{"_id": userId}}}
	// Find the user document in the "organizations" collection
	result, err := r.mongodb.Collection("organizations").UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// Check if the update operation modified any documents
	if result.ModifiedCount == 0 {
		return err
	}

	return nil
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
