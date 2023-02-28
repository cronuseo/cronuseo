package monitoring

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/entity"
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

	data := entity.AllowedData{}

	return data, nil
}
