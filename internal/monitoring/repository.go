package monitoring

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetAllowed(ctx context.Context, orgId string) (entity.AllowedData, error)
}

type repository struct {
	mongodb *mongo.Database
}

func NewRepository(mongodb *mongo.Database) Repository {

	return repository{mongodb: mongodb}
}

// Get allowed data.
func (r repository) GetAllowed(ctx context.Context, orgId string) (entity.AllowedData, error) {

	collection := r.mongodb.Collection("checks")

	data := entity.AllowedData{}

	trueCount, err := collection.CountDocuments(context.Background(), bson.M{"result": true})
	if err != nil {
		return data, err
	}

	falseCount, err := collection.CountDocuments(context.Background(), bson.M{"result": false})
	if err != nil {
		return data, err
	}
	data.Allowed = int(trueCount)
	data.NotAllowed = int(falseCount)
	return data, nil
}
