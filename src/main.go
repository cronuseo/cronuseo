package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/routes"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	router := gin.Default()
	config.ConnectDB()
	config.InitLogger()
	routes.OAuth2Routes(router)
	routes.UserRoutes(router)
	routes.OrganizationRoutes(router)
	routes.ProjectRoutes(router)
	routes.ResourceRoutes(router)
	routes.ResourceActionRoutes(router)
	routes.ResourceRoleRoutes(router)
	router.Run() // listen and serve on 0.0.0.0:8080
}
