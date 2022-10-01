package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func OrganizationRoutes(router *gin.Engine) {

	router.GET("/orgs", controllers.GetOrganizations)
	router.POST("/orgs", controllers.CreateOrganization)
	router.DELETE("/orgs/:id", controllers.DeleteOrganization)
	router.POST("/orgs/:id", controllers.UpdateOrganization)
}
