package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/shashimalcse/cronuseo/internal/config"
	"github.com/shashimalcse/cronuseo/internal/keto"
	"github.com/shashimalcse/cronuseo/internal/organization"
	"github.com/shashimalcse/cronuseo/internal/permission"
	"github.com/shashimalcse/cronuseo/internal/resource"
	"github.com/shashimalcse/cronuseo/internal/role"
	"github.com/shashimalcse/cronuseo/internal/user"
	"google.golang.org/grpc"

	_ "github.com/shashimalcse/cronuseo/docs"

	jwt "github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	jwk "github.com/lestrrat-go/jwx/jwk"
	_ "github.com/lib/pq"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

// var flagConfig = flag.String("config", "/Users/thilinashashimal/Desktop/Cronuseo/config/local.yml", "path to the config file")

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
		log.Fatalln("error connecting databse")
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

	clients := keto.KetoClients{WriteClient: writeClient, ReadClient: readClient, CheckClient: checkClient}

	e := buildHandler(db, cfg, clients)
	e.Logger.Fatal(e.Start(cfg.API))

}

func buildHandler(db *sqlx.DB, cfg *config.Config, clients keto.KetoClients) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORS())
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	rg := router.Group("/api/v1")
	// config := middleware.JWTConfig{
	// 	KeyFunc: getKey(cfg),
	// }
	// rg.Use(middleware.JWTWithConfig(config))
	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(db)))
	user.RegisterHandlers(rg, user.NewService(user.NewRepository(db)))
	resource.RegisterHandlers(rg, resource.NewService(resource.NewRepository(db)))
	role.RegisterHandlers(rg, role.NewService(role.NewRepository(db)))
	permission.RegisterHandlers(rg, permission.NewService(permission.NewRepository(db)))
	keto.RegisterHandlers(rg, keto.NewService(keto.NewRepository(clients, db)))
	return router
}

func getKey(cfg *config.Config) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {

		keySet, err := jwk.Fetch(context.Background(), cfg.JWKS)
		if err != nil {
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have a key ID in the kid field")
		}

		key, found := keySet.LookupKeyID(keyID)

		if !found {
			return nil, fmt.Errorf("unable to find key %q", keyID)
		}

		var pubkey interface{}
		if err := key.Raw(&pubkey); err != nil {
			return nil, fmt.Errorf("Unable to get the public key. Error: %s", err.Error())
		}

		return pubkey, nil
	}
}
