package main

import (
	"cronuseo/internal/config"
	"cronuseo/internal/organization"
	"flag"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

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

	handler := buildHandler(db, cfg)
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	handler.Logger.Fatal(handler.Start(address))

}

func buildHandler(db *sqlx.DB, cfg *config.Config) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORS())

	rg := router.Group("/api/v1")

	organization.RegisterHandlers(rg, organization.NewService(organization.NewRepository(db)))
	return router
}
