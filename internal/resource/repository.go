package resource

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
	Get(ctx context.Context, org_id string, id string) (*mongo_entity.Resource, error)
	// Query(ctx context.Context, org_id string, filter Filter) ([]mongo_entity.Resource, error)
	Create(ctx context.Context, org_id string, resource mongo_entity.Resource) error
	// Update(ctx context.Context, org_id string, resource mongo_entity.Resource) error
	// Delete(ctx context.Context, org_id string, id string) error
	CheckResourceExistById(ctx context.Context, org_id string, id string) (bool, error)
	CheckResourceExistsByIdentifier(ctx context.Context, org_id string, key string) (bool, error)
}

type repository struct {
	mongodb *mongo.Database
}

func NewRepository(mongodb *mongo.Database) Repository {

	return repository{mongodb: mongodb}
}

// Get resource by id.
func (r repository) Get(ctx context.Context, org_id string, id string) (*mongo_entity.Resource, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return nil, err
	}

	resId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define filter to find the resource by its ID
	filter := bson.M{"_id": orgId, "resources._id": resId}
	projection := bson.M{"resources.$": 1}
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

	return &org.Resources[0], nil
}

// Create new resource.
func (r repository) Create(ctx context.Context, org_id string, resource mongo_entity.Resource) error {

	coll := r.mongodb.Collection("organizations")

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return err
	}
	// Update the APIResources array for the given organization
	filter := bson.M{"_id": orgId}
	update := bson.M{"$push": bson.M{"resources": resource}}
	_, err = coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil

}

// Update resource.
// func (r repository) Update(ctx context.Context, org_id string, resource mongo_entity.Resource) error {

// 	stmt, err := r.db.Prepare("UPDATE org_resource SET name = $1 HERE org_id = $2 AND resource_id = $3")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = stmt.Exec(resource.Name, org_id, resource.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Delete resource.
// func (r repository) Delete(ctx context.Context, org_id string, id string) error {

// 	tx, err := r.db.DB.Begin()

// 	if err != nil {
// 		return err
// 	}

// 	// Delete actions assigned to the resource
// 	{
// 		stmt, err := tx.Prepare("DELETE FROM res_action WHERE resource_id = $1")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		_, err = stmt.Exec(id)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Delete resource
// 	{
// 		stmt, err := tx.Prepare("DELETE FROM org_resource WHERE org_id = $1 AND resource_id = $2")
// 		if err != nil {
// 			return err
// 		}
// 		defer stmt.Close()
// 		_, err = stmt.Exec(org_id, id)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	{
// 		err := tx.Commit()

// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// Get all resources.
// func (r repository) Query(ctx context.Context, org_id string, filter Filter) ([]entity.Resource, error) {

// 	resources := []entity.Resource{}
// 	name := filter.Name + "%"
// 	err := r.db.Select(&resources, "SELECT * FROM org_resource WHERE org_id = $1 AND name LIKE $2 AND id > $3 ORDER BY id LIMIT $4", org_id, name, filter.Cursor, filter.Limit)
// 	return resources, err
// }

// Check if resource exists by id.
func (r repository) CheckResourceExistById(ctx context.Context, org_id string, id string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}

	resId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": orgId, "resources._id": resId}

	// Search for the resource in the "organizations" collection
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

// Check if resource exists by key.
func (r repository) CheckResourceExistsByIdentifier(ctx context.Context, org_id string, identifier string) (bool, error) {

	orgId, err := primitive.ObjectIDFromHex(org_id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": orgId, "resources.identifier": identifier}

	// Search for the resource in the "organizations" collection
	count, err := r.mongodb.Collection("organizations").CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
