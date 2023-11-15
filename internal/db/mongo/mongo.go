package mongo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/bson"
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
	// Initialize default organization
	initializeOrganization(mongodb, cfg, logger)
	return mongodb, nil
}

func initializeOrganization(mongodb *MongoDB, cfg *config.Config, logger *zap.Logger) {

	orgCollection := mongodb.MongoClient.Database(mongodb.MongoConfig.DBName).Collection(mongodb.MongoConfig.OrganizationCollectionName)
	filter := bson.M{"identifier": cfg.RootOrganization.Name}
	var org mongo_entity.Organization
	err := orgCollection.FindOne(context.Background(), filter).Decode(&org)
	if err == mongo.ErrNoDocuments {

		// Organization doesn't exist, so create it
		key := make([]byte, 32)

		if _, err := rand.Read(key); err != nil {
			logger.Fatal("Error while initializing organization", zap.String("identifier", cfg.RootOrganization.Name), zap.String("error", err.Error()))
			os.Exit(-1)

		}

		APIKey := base64.StdEncoding.EncodeToString(key)

		defaultOrg := mongo_entity.Organization{
			DisplayName:     cfg.RootOrganization.Name,
			Identifier:      cfg.RootOrganization.Name,
			API_KEY:         APIKey,
			Resources:       []mongo_entity.Resource{},
			Users:           []mongo_entity.User{},
			Roles:           []mongo_entity.Role{},
			Groups:          []mongo_entity.Group{},
			RolePermissions: []mongo_entity.RolePermission{},
		}
		_, err = orgCollection.InsertOne(context.Background(), defaultOrg)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("Default organization created")
	} else if err != nil {
		logger.Fatal("Error while initializing organization", zap.String("identifier", cfg.RootOrganization.Name), zap.String("error", err.Error()))
	} else {
		logger.Info("Organization already exists!", zap.String("identifier", cfg.RootOrganization.Name))
	}
}
