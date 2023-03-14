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
	monitoringClient *mongo.Client
}

func NewRepository(monitoringClient *mongo.Client) Repository {

	return repository{monitoringClient: monitoringClient}
}

// Get allowed data.
func (r repository) GetAllowed(ctx context.Context, orgId string) (entity.AllowedData, error) {

	collection := r.monitoringClient.Database("monitoring").Collection("checks")

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
