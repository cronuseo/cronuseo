package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	routes.OrganizationRoutes2(e)
	routes.ProjectRoutes2(e)
	routes.ResourceRoutes2(e)

	e.Logger.Fatal(e.Start(":8080"))
}
