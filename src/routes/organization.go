package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func OrganizationRoutes(router *gin.Engine) {

	orgRouter := router.Group("/orgs")

	orgRouter.GET("/", controllers.GetOrganizations)
	orgRouter.GET("/:id", controllers.GetOrganization)
	orgRouter.POST("/", controllers.CreateOrganization)
	orgRouter.DELETE("/:id", controllers.DeleteOrganization)
	orgRouter.PUT("/:id", controllers.UpdateOrganization)
}
