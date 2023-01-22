package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/shashimalcse/cronuseo/internal/action"
	"github.com/shashimalcse/cronuseo/internal/auth"
	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/check"
	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/permission"
	"github.com/shashimalcse/cronuseo/internal/resource"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/user"
	"github.com/shashimalcse/cronuseo/internal/util"
	"google.golang.org/grpc"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	_ "github.com/shashimalcse/cronuseo/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	flag.Parse()

	// Load configurations for db, keto and redis.
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatal("error loading config")
		os.Exit(-1)
	}

	// Connect to database.
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		log.Fatalln("error connecting database")
		os.Exit(-1)
	}

	// We need three keto client for checking, reading, and writing.

	// Write client.
	conn, err := grpc.Dial(cfg.KetoWrite, grpc.WithInsecure())
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
	writeClient := rts.NewWriteServiceClient(conn)

	// Read client.
	conn, err = grpc.Dial(cfg.KetoRead, grpc.WithInsecure())
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
	readClient := rts.NewReadServiceClient(conn)

	// Check client.
	conn, err = grpc.Dial(cfg.KetoRead, grpc.WithInsecure())
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
	checkClient := rts.NewCheckServiceClient(conn)

	// Create a struct to hold all keto clients.
	clients := util.KetoClients{
		WriteClient: writeClient,
		ReadClient:  readClient,
		CheckClient: checkClient}

	// Redis client.
	permissionCache := cache.NewRedisCache("localhost:6379", 0, 200)

	e := buildHandler(db, cfg, clients, permissionCache)
	e.Logger.Fatal(e.Start(cfg.API))

}

// buildHandler builds the echo router.
func buildHandler(db *sqlx.DB, cfg *config.Config, clients util.KetoClients, permissionCache cache.PermissionCache) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "API_KEY"}, // API_KEY is used for permission checking SDKs
		AllowOrigins:     []string{"http://localhost:3000"},
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
	check.RegisterHandlers(rg, check.NewService(check.NewRepository(clients, db), permissionCache))
	permission.RegisterHandlers(rg, permission.NewService(permission.NewRepository(clients, db), permissionCache))
	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(db)))
	user.RegisterHandlers(rg, user.NewService(user.NewRepository(db), permissionCache))
	resource.RegisterHandlers(rg, resource.NewService(resource.NewRepository(db)))
	role.RegisterHandlers(rg, role.NewService(role.NewRepository(db, clients.WriteClient), permissionCache))
	action.RegisterHandlers(rg, action.NewService(action.NewRepository(db)))

	return router
}
