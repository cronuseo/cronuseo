package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/routes"
	"github.com/shashimalcse/Cronuseo/utils"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	config.ConnectDB()
	config.InitLogger()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	routes.OrganizationRoutes(e)
	routes.ProjectRoutes(e)
	routes.ResourceRoutes(e)
	routes.ResourceActionRoutes(e)
	routes.ResourceRoleRoutes(e)
	routes.ResourceRoutes(e)
	routes.GroupRoutes(e)
	routes.UserRoutes(e)
	routes.CheckRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
