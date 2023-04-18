package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	_ "github.com/shashimalcse/cronuseo/docs"
	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/group"
	"github.com/shashimalcse/cronuseo/internal/mongo_entity"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/resource"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/user"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
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

	var asgardeo = oauth2.Endpoint{
		AuthURL:  "https://api.asgardeo.io/t/cronuseo/oauth2/authorize",
		TokenURL: "https://api.asgardeo.io/t/cronuseo/oauth2/token",
	}

	asgardeoOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/v1/auth/callback",
		ClientID:     "PrFxAZefrbkVPN6RkLATzUAlbUga",
		ClientSecret: "tzMbQz83mx7HPXxIej2uKIcQtoAa",
		Scopes:       []string{"openid", "profile"},
		Endpoint:     asgardeo,
	}

	e := buildHandler(cfg, logger, mongo_db, asgardeoOauthConfig)
	logger.Info("Starting server", zap.String("server_endpoint", cfg.Mgt_API))
	e.Logger.Fatal(e.Start(cfg.Mgt_API))

}

// buildHandler builds the echo router.
func buildHandler(
	cfg *config.Config, // Config
	logger *zap.Logger, // Logger
	mongodb *mongo.Database, // Mongo monitoring client
	asgardeoOauthConfig *oauth2.Config,
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
	rg.Use(authmw(cfg.JWKS))

	// Here we register all the handlers. Each handler handle jwt validation separately.
	// check.RegisterHandlers(rg, check.NewService(check.NewRepository(clients, db), permissionCache, logger), mongodb.Collection("checks"))
	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(mongodb), logger))
	user.RegisterHandlers(rg, user.NewService(user.NewRepository(mongodb), logger))
	resource.RegisterHandlers(rg, resource.NewService(resource.NewRepository(mongodb), logger))
	role.RegisterHandlers(rg, role.NewService(role.NewRepository(mongodb), logger))
	group.RegisterHandlers(rg, group.NewService(group.NewRepository(mongodb), logger))

	return router
}

func authmw(jwksURL string) echo.MiddlewareFunc {

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			panic(fmt.Sprintf("There was an error with the jwt.KeyFunc\nError:%s\n", err.Error()))
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to create JWKs from resource at the given URL.\nError:%s\n", err.Error()))
	}

	// initialize JWT middleware instance
	return middleware.JWTWithConfig(middleware.JWTConfig{
		// skip public endpoints
		// Skipper: func(context echo.Context) bool {
		// 	return context.Path() == "/"
		// },
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			// a hack to convert jwt -> v4 tokens
			t, _, err := new(jwtv4.Parser).ParseUnverified(token.Raw, jwtv4.MapClaims{})
			if err != nil {
				return nil, err
			}
			return jwks.Keyfunc(t)
		},
	})
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
			DisplayName:     orgName,
			Identifier:      orgName,
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
		logger.Fatal("Error while initializing organization", zap.String("initializing_organization", orgName))
	} else {
		logger.Info("Organization already exists")
	}
}
