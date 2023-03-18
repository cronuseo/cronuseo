package organization

import (
	"context"
	"fmt"

	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, id string) (*mongo_entity.Organization, error)
	Query(ctx context.Context) ([]mongo_entity.Organization, error)
	Create(ctx context.Context, organization mongo_entity.Organization) error
	Delete(ctx context.Context, id string) error
	RefreshAPIKey(ctx context.Context, apiKey string, id string) error
	CheckOrgExistById(ctx context.Context, id string) (bool, error)
	CheckOrgExistByIdentifier(ctx context.Context, identifier string) (bool, error)
}

type repository struct {
	db      *sqlx.DB
	mongodb *mongo.Database
}

func NewRepository(db *sqlx.DB, mongodb *mongo.Database) Repository {

	return repository{db: db, mongodb: mongodb}
}

// Get organization by id.
func (r repository) Get(ctx context.Context, id string) (*mongo_entity.Organization, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the organization by its ID
	filter := bson.M{"_id": objID}

	// Find the organization document in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)
	if err := result.Err(); err != nil {
		return nil, err
	}

	// Decode the organization document into a struct
	var org mongo_entity.Organization
	if err := result.Decode(&org); err != nil {
		return nil, err
	}

	return &org, nil
}

// Create new organization.
func (r repository) Create(ctx context.Context, organization mongo_entity.Organization) error {

	_, err := r.mongodb.Collection("organizations").InsertOne(context.Background(), organization)
	if err != nil {
		return err
	}
	return nil
}

// Delete organization.
func (r repository) Delete(ctx context.Context, id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define the filter to search for the organization by ID
	filter := bson.M{"_id": objID}

	// Delete the organization from the "organizations" collection
	result, err := r.mongodb.Collection("organizations").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	// Check if the organization was deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("Organization with ID %s not found", id)
	}

	return nil
}

// Refresh API key in mongo.
func (r repository) RefreshAPIKey(ctx context.Context, apiKey string, id string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Define filter to find the organization by its ID
	filter := bson.M{"_id": objID}

	// Define update to set the new name
	update := bson.M{"$set": bson.M{"api_key": apiKey}}

	// Define options for update operation
	options := options.Update().SetUpsert(false)

	// Update the organization document in the "organizations" collection
	result, err := r.mongodb.Collection("organizations").UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return err
	}

	// Check if any documents were updated
	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Query organizations.
func (r repository) Query(ctx context.Context) ([]mongo_entity.Organization, error) {

	// Define an empty slice to store the organizations
	var orgs []mongo_entity.Organization

	// Search for all organizations in the "organizations" collection
	cursor, err := r.mongodb.Collection("organizations").Find(context.Background(), bson.M{})
	if err != nil {
		return orgs, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the results and add each organization to the slice
	for cursor.Next(context.Background()) {
		var org mongo_entity.Organization
		if err := cursor.Decode(&org); err != nil {
			return orgs, err
		}
		orgs = append(orgs, org)
	}

	// Check for any errors that occurred during iteration
	if err := cursor.Err(); err != nil {
		return orgs, err
	}

	return orgs, nil
}

// Check if organization exists by id.
func (r repository) CheckOrgExistById(ctx context.Context, id string) (bool, error) {

	filter := bson.M{"_id": id}

	// Search for the organization in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)

	// Check if the organization was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}

// Check if organization exists by identifier.
func (r repository) CheckOrgExistByIdentifier(ctx context.Context, identifier string) (bool, error) {

	filter := bson.M{"identifier": identifier}

	// Search for the organization in the "organizations" collection
	result := r.mongodb.Collection("organizations").FindOne(context.Background(), filter)

	// Check if the organization was found
	if result.Err() == nil {
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else {
		return false, result.Err()
	}
}
