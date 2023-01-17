package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/shashimalcse/cronuseo/internal/action"
	"github.com/shashimalcse/cronuseo/internal/auth"
	"github.com/shashimalcse/cronuseo/internal/cache"
	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/permission"
	"github.com/shashimalcse/cronuseo/internal/resource"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/user"
	"google.golang.org/grpc"

	"github.com/jmoiron/sqlx"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	_ "github.com/shashimalcse/cronuseo/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var Version = "1.0.0"

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
	// load application configurations
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatal("error loading config")
		os.Exit(-1)
	}

	//connect db
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		log.Fatalln("error connecting database")
		os.Exit(-1)
	}

	//keto clients

	conn, err := grpc.Dial(cfg.KetoWrite, grpc.WithInsecure())
	if err != nil {
		panic("Encountered error: " + err.Error())
	}

	writeClient := rts.NewWriteServiceClient(conn)

	conn, err = grpc.Dial(cfg.KetoRead, grpc.WithInsecure())
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
	readClient := rts.NewReadServiceClient(conn)

	conn, err = grpc.Dial(cfg.KetoRead, grpc.WithInsecure())
	if err != nil {
		panic("Encountered error: " + err.Error())
	}
	checkClient := rts.NewCheckServiceClient(conn)

	clients := permission.KetoClients{WriteClient: writeClient, ReadClient: readClient, CheckClient: checkClient}

	// redis client
	permissionCache := cache.NewRedisCache("localhost:6379", 0, 200)

	e := buildHandler(db, cfg, clients, permissionCache)
	e.Logger.Fatal(e.Start(cfg.API))

}

func buildHandler(db *sqlx.DB, cfg *config.Config, clients permission.KetoClients, permissionCache cache.PermissionCache) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowCredentials},
	}))
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	rg := router.Group("/api/v1")
	rg.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	auth.RegisterHandlers(rg, auth.NewService(auth.NewRepository(db)), organization.NewService(organization.NewRepository(db)))
	config := echojwt.Config{
		SigningKey: []byte(auth.SecretKey),
	}
	rg.Use(echojwt.WithConfig(config))
	// Define the health endpoint
	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(db)))
	user.RegisterHandlers(rg, user.NewService(user.NewRepository(db)))
	resource.RegisterHandlers(rg, resource.NewService(resource.NewRepository(db)))
	role.RegisterHandlers(rg, role.NewService(role.NewRepository(db)))
	action.RegisterHandlers(rg, action.NewService(action.NewRepository(db)))
	permission.RegisterHandlers(rg, permission.NewService(permission.NewRepository(clients, db), permissionCache))
	return router
}
