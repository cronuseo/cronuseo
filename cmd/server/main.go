package main

import (
	"cronuseo/internal/config"
	"cronuseo/internal/organization"
	"flag"
	"os"

	_ "cronuseo/docs"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
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
		os.Exit(-1)
	}

	//connect db
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		os.Exit(-1)
	}
	print("start server")
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.CORS())
	apiV1 := e.Group("/api/v1")
	organization.RegisterHandlers(apiV1, organization.NewService(organization.NewRepository(db)))
	e.Logger.Fatal(e.Start(":8080"))

	// address := fmt.Sprintf(":%v", cfg.ServerPort)
	// e.Logger.Fatal(e.Start(address))

}

func buildHandler(db *sqlx.DB, cfg *config.Config) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORS())
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	rg := router.Group("/api/v1")

	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(db)))
	return router
}
