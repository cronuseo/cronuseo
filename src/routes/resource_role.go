package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoleRoutes(router *echo.Group) {

	resourceRoleRouter := router.Group("/:res_id/resource_role")

	resourceRoleRouter.GET("/:id", controllers.GetResourceRole)
	resourceRoleRouter.GET("", controllers.GetResourceRoles)
	resourceRoleRouter.POST("", controllers.CreateResourceRole)
	resourceRoleRouter.PATCH("/:id", controllers.PatchResourceRole)
	resourceRoleRouter.DELETE("/:id", controllers.DeleteResourceRole)
	resourceRoleRouter.PUT("/:id", controllers.UpdateResourceRole)
}
