package main

import (
	"flag"
	"log"
	"os"

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
		log.Fatal("Error while loading config.", zap.String("config_file", flag.Lookup("config").Value.String()))
		os.Exit(-1)
	}

	// Set up logger.
	logger := logger.Init(cfg)

	// Mongo client.
	mongodb := db.Init(cfg, logger)

	e := buildHandler(cfg, logger, mongodb)
	logger.Info("Starting server", zap.String("server_endpoint", cfg.Mgt_API))
	e.Logger.Fatal(e.Start(cfg.Mgt_API))

}

// buildHandler builds the echo router.
func buildHandler(
	cfg *config.Config, // Config
	logger *zap.Logger, // Logger
	mongodb *db.MongoDB, // MongoDB
) *echo.Echo {

	router := echo.New()

	// Set up CORS.
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
	rg.Use(mw.Auth(cfg, logger))

	// Here we register all the handlers.
	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(mongodb), logger))
	user.RegisterHandlers(rg, user.NewService(user.NewRepository(mongodb), logger))
	resource.RegisterHandlers(rg, resource.NewService(resource.NewRepository(mongodb), logger))
	role.RegisterHandlers(rg, role.NewService(role.NewRepository(mongodb), logger))
	group.RegisterHandlers(rg, group.NewService(group.NewRepository(mongodb), logger))

	return router
}
