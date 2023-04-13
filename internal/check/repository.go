package check

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
}

type repo struct {
	mongodb *mongo.Database
}

func NewRepository(mongodb *mongo.Database) Repository {
	return repo{mongodb: mongodb}
}
