package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func OrganizationRoutes(router *echo.Group) {

	orgRouter := router.Group("/organization")

	orgRouter.GET("", controllers.GetOrganizations)
	orgRouter.GET("/:id", controllers.GetOrganization)
	orgRouter.POST("", controllers.CreateOrganization)
	orgRouter.DELETE("/:id", controllers.DeleteOrganization)
	orgRouter.PUT("/:id", controllers.UpdateOrganization)
}
