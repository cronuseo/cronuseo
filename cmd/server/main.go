package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	_ "github.com/shashimalcse/cronuseo/docs"
	"github.com/shashimalcse/cronuseo/internal/action"
	"github.com/shashimalcse/cronuseo/internal/auth"
	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/check"
	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/monitoring"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/permission"
	"github.com/shashimalcse/cronuseo/internal/resource"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/user"
	"github.com/shashimalcse/cronuseo/internal/util"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var Version = "1.0.0"

// Default config flag.
var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

// @title          Cronuseo API
// @version        1.0
// @description    This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath /api/v1
func main() {

	logger := InitializeLogger()

	flag.Parse()

	// Load configurations for db, keto and redis.
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		logger.Fatal("Error while loading config.", zap.String("config_file", flag.Lookup("config").Value.String()))
		os.Exit(-1)
	}

	// Connect to database.
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		logger.Fatal("Error while connecting to the database.", zap.String("database_url", cfg.DSN))
		os.Exit(-1)
	}

	// We need three keto client for checking, reading, and writing.

	// Write client.
	conn, err := grpc.Dial(cfg.KetoWrite, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Error while creating keto write client", zap.String("write_client_endpoint", cfg.KetoWrite))
		os.Exit(-1)
	}
	writeClient := rts.NewWriteServiceClient(conn)

	// Read client.
	conn, err = grpc.Dial(cfg.KetoRead, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Error while creating keto read client", zap.String("read_client_endpoint", cfg.KetoWrite))
		os.Exit(-1)
	}
	readClient := rts.NewReadServiceClient(conn)

	// Check client.
	conn, err = grpc.Dial(cfg.KetoRead, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Error while creating keto check client", zap.String("check_client_endpoint", cfg.KetoWrite))
		os.Exit(-1)
	}
	checkClient := rts.NewCheckServiceClient(conn)

	// Create a struct to hold all keto clients.
	clients := util.KetoClients{
		WriteClient: writeClient,
		ReadClient:  readClient,
		CheckClient: checkClient,
	}

	// Redis client.
	permissionCache := cache.NewRedisCache(cfg.RedisEndpoint, 0, 200, cfg.RedisPassword)

	// Mongo client.
	credential := options.Credential{
		Username: cfg.MongoUser,
		Password: cfg.MongoPassword,
	}
	mongo_client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Mongo).SetAuth(credential))
	if err != nil {
		panic(err)
	}

	mongo_db := mongo_client.Database(cfg.MongoDBName)
	InitializeOrganization(mongo_db, logger, cfg.DefaultOrg)

	e := buildHandler(db, cfg, logger, clients, permissionCache, mongo_db)
	logger.Info("Starting server", zap.String("server_endpoint", cfg.API))
	e.Logger.Fatal(e.Start(cfg.API))

}

// buildHandler builds the echo router.
func buildHandler(
	db *sqlx.DB, // Database
	cfg *config.Config, // Config
	logger *zap.Logger, // Logger
	clients util.KetoClients, // Keto clients
	permissionCache cache.PermissionCache, // Redis permission cache
	mongodb *mongo.Database, // Mongo monitoring client
) *echo.Echo {

	router := echo.New()

	//  Set up CORS.
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "API_KEY"}, // API_KEY is used for permission checking SDKs
		AllowOrigins:     []string{"http://localhost:3000"},
	}))
	// Echo logger middleware.
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}; method=${method}; uri=${uri}; status=${status};\n",
	}))

	// Swagger endpoint.
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	// API endpoints.
	rg := router.Group("/api/v1")

	// Currently this health check endpoint is used by the run console script to check availability of the service.
	rg.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Here we register all the handlers. Each handler handle jwt validation separately.
	auth.RegisterHandlers(rg, auth.NewService(auth.NewRepository(db)))
	check.RegisterHandlers(rg, check.NewService(check.NewRepository(clients, db), permissionCache, logger), mongodb.Collection("checks"))
	permission.RegisterHandlers(rg, permission.NewService(permission.NewRepository(clients, db), permissionCache, logger))
	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(mongodb), logger, permissionCache))
	user.RegisterHandlers(rg, user.NewService(user.NewRepository(mongodb), permissionCache, logger))
	resource.RegisterHandlers(rg, resource.NewService(resource.NewRepository(mongodb), logger))
	role.RegisterHandlers(rg, role.NewService(role.NewRepository(db, clients.WriteClient), permissionCache, logger))
	action.RegisterHandlers(rg, action.NewService(action.NewRepository(db), logger))
	monitoring.RegisterHandlers(rg, monitoring.NewService(monitoring.NewRepository(mongodb), logger))

	return router
}

func InitializeLogger() *zap.Logger {

	zap_config := zap.NewProductionEncoderConfig()
	zap_config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zap_config)
	consoleEncoder := zapcore.NewConsoleEncoder(zap_config)
	logFile, _ := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func InitializeOrganization(db *mongo.Database, logger *zap.Logger, orgName string) {

	orgCollection := db.Collection("organizations")
	filter := bson.M{"identifier": orgName}
	var org mongo_entity.Organization
	err := orgCollection.FindOne(context.Background(), filter).Decode(&org)
	if err == mongo.ErrNoDocuments {
		// Organization doesn't exist, so create it

		key := make([]byte, 32)

		if _, err := rand.Read(key); err != nil {
			logger.Fatal("Error while initializing organization", zap.String("initializing_organization", orgName))
			os.Exit(-1)

		}
		APIKey := base64.StdEncoding.EncodeToString(key)

		defaultOrg := mongo_entity.Organization{
			DisplayName: orgName,
			Identifier:  orgName,
			API_KEY:     APIKey,
			Resources:   []mongo_entity.Resource{},
			Users:       []mongo_entity.User{},
			Roles:       []mongo_entity.Role{},
		}
		_, err = orgCollection.InsertOne(context.Background(), defaultOrg)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("Default organization created")
	} else if err != nil {
		logger.Fatal("Error while initializing organization", zap.String("initializing_organization", orgName))
	} else {
		logger.Info("Organization already exists")
	}
}
