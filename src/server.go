package main

import (
	"log"

	_ "github.com/shashimalcse/Cronuseo/docs"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/routes"
	"github.com/shashimalcse/Cronuseo/utils"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//  @title           Cronuseo API
//  @version         1.0
//  @description     This is a sample server celler server.
//  @termsOfService  http://swagger.io/terms/

//  @contact.name   API Support
//  @contact.url    http://www.swagger.io/support
//  @contact.email  support@swagger.io

//  @license.name  Apache 2.0
//  @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	config.ConnectDB()
	config.InitLogger()
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.CORS())
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	setRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func setRoutes(e *echo.Echo) {
	apiV1 := e.Group("/api/v1")
	routes.OrganizationRoutes(apiV1)
	routes.ProjectRoutes(apiV1)
	routes.ResourceRoutes(apiV1)
	routes.ResourceActionRoutes(apiV1)
	routes.ResourceRoleRoutes(apiV1)
	routes.ResourceRoutes(apiV1)
	routes.GroupRoutes(apiV1)
	routes.UserRoutes(apiV1)
	routes.CheckRoutes(apiV1)
}
