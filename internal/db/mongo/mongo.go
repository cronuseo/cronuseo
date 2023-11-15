package mongo

import (
	"context"

	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	MongoClient *mongo.Client
	MongoConfig util.MongoDBConfig
}

func Init(cfg *config.Config, logger *zap.Logger) (*MongoDB, error) {

	credential := options.Credential{
		Username: cfg.Database.User,
		Password: cfg.Database.Password,
	}
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Database.URL).SetAuth(credential))
	if err != nil {
		logger.Fatal("Error while connecting to MongoDB", zap.String("error", err.Error()))
		return nil, err
	}

	mongoConfig := util.MongoDBConfig{
		DBName:                     cfg.Database.Name,
		OrganizationCollectionName: cfg.Database.Name,
	}

	mongodb := &MongoDB{MongoClient: mongoClient, MongoConfig: mongoConfig}
	return mongodb, nil
}
