package main

import (
	"flag"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	_ "github.com/shashimalcse/cronuseo/docs"
	"github.com/shashimalcse/cronuseo/internal/config"
	db "github.com/shashimalcse/cronuseo/internal/db/mongo"
	"github.com/shashimalcse/cronuseo/internal/group"
	"github.com/shashimalcse/cronuseo/internal/logger"
	mw "github.com/shashimalcse/cronuseo/internal/middleware"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/resource"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/user"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var Version = "1.0.0"

// Default config flag.
var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

// @title          Cronuseo API
// @version        1.0
// @description    This is a cronuseo server.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath /api/v1
func main() {

	flag.Parse()

	// Load configurations.
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatalf("Error while loading config: %v\n", err)
	}

	// Set up logger.
	logger, err := logger.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v\n", err)
	}

	// Mongo client.
	mongodb, err := db.Init(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB client", zap.Error(err))
	}

	logger.Info("Starting server", zap.String("server_endpoint", cfg.Endpoint.Management))

	if err := BuildServer(cfg, logger, mongodb).Start(cfg.Endpoint.Management); err != nil {
		logger.Fatal("Error while starting server", zap.Error(err))
	}
}

// BuildServer builds and configures the echo server.
func BuildServer(
	cfg *config.Config, // Config
	logger *zap.Logger, // Logger
	mongodb *db.MongoDB, // MongoDB
) *echo.Echo {

	e := echo.New()

	// Middleware setup.
	setupMiddleware(e, cfg)

	// API route groups.
	apiV1 := e.Group("/api/v1")

	// Apply middleware specific to API routes if needed.
	apiV1.Use(mw.Auth(cfg, logger))

	// Register service handlers.
	registerServiceHandlers(apiV1, mongodb, cfg, logger)

	return e
}

func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	// CORS middleware configuration.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "API_KEY"},
		AllowOrigins:     []string{"http://localhost:3000"},
	}))

	// Logger middleware.
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}; method=${method}; uri=${uri}; status=${status};\n",
	}))

	// Swagger endpoint.
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

func registerServiceHandlers(e *echo.Group, mongodb *db.MongoDB, cfg *config.Config, logger *zap.Logger) {
	// Initialize repositories.
	orgRepo := organization.NewRepository(mongodb)
	userRepo := user.NewRepository(mongodb)
	resourceRepo := resource.NewRepository(mongodb)
	roleRepo := role.NewRepository(mongodb)
	groupRepo := group.NewRepository(mongodb)

	// Initialize services with repositories.
	orgService := organization.NewService(orgRepo, logger)
	userService := user.NewService(userRepo, logger)
	resourceService := resource.NewService(resourceRepo, logger)
	roleService := role.NewService(roleRepo, logger)
	groupService := group.NewService(groupRepo, logger)

	initializeAdmin(orgService, userService, roleService, cfg, logger)

	// Register handlers.
	organization.RegisterHandlers(e, orgService)
	user.RegisterHandlers(e, userService)
	resource.RegisterHandlers(e, resourceService)
	role.RegisterHandlers(e, roleService)
	group.RegisterHandlers(e, groupService)
}

func initializeAdmin(orgService organization.Service, userService user.Service, roleService role.Service, cfg *config.Config, logger *zap.Logger) {

	rootOrgId, err := orgService.GetIdByIdentifier(nil, cfg.RootOrganization.Name)
	if err != nil {
		logger.Fatal("Failed to get root org id", zap.Error(err))
	}
	adminUser := user.CreateUserRequest{
		Identifier: cfg.RootOrganization.AdminIdentifier,
		Username:   cfg.RootOrganization.AdminName,
		Roles:      []primitive.ObjectID{},
		Groups:     []primitive.ObjectID{},
	}
	userService.Create(nil, rootOrgId, adminUser)
	adminId, err := userService.GetIdByIdentifier(nil, rootOrgId, cfg.RootOrganization.AdminIdentifier)
	if err != nil {
		logger.Fatal("Failed to get admin id", zap.Error(err))
	}
	adminObjID, _ := primitive.ObjectIDFromHex(adminId)
	adminRole := role.CreateRoleRequest{
		Identifier:  cfg.RootOrganization.AdminRoleName,
		DisplayName: cfg.RootOrganization.AdminRoleName,
		Users: []primitive.ObjectID{
			adminObjID,
		},
		Groups: []primitive.ObjectID{},
	}
	roleService.Create(nil, rootOrgId, adminRole)
}
