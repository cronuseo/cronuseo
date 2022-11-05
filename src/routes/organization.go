package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func OrganizationRoutes(router *echo.Echo) {

	orgRouter := router.Group("/orgs")

	orgRouter.GET("", controllers.GetOrganizations)
	orgRouter.GET("/:id", controllers.GetOrganization)
	orgRouter.POST("", controllers.CreateOrganization)
	orgRouter.DELETE("/:id", controllers.DeleteOrganization)
	orgRouter.PUT("/:id", controllers.UpdateOrganization)
}
