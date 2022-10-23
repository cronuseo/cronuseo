package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
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

func OrganizationRoutes2(router *echo.Echo) {

	orgRouter := router.Group("/organization")

	orgRouter.GET("", controllers.GetOrganizations2)
	orgRouter.GET("/:id", controllers.GetOrganization2)
	orgRouter.POST("", controllers.CreateOrganization2)
	orgRouter.DELETE("/:id", controllers.DeleteOrganization2)
	orgRouter.PUT("/:id", controllers.UpdateOrganization2)

}
