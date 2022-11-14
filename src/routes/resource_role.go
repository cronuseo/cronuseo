package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoleRoutes(router *echo.Group) {

	resourceActionRouter := router.Group("/:res_id/resource_role")

	resourceActionRouter.GET("/:id", controllers.GetResourceRole)
	resourceActionRouter.GET("", controllers.GetResourceRoles)
	resourceActionRouter.POST("", controllers.CreateResourceRole)
	// resourceActionRouter.POST("/:id/user/:user_id", controllers.AddUserToResourceRole)
	// resourceActionRouter.POST("/:id/group/:group_id", controllers.AddGroupToResourceRole)
	// resourceActionRouter.POST("/:id/action/:resact_id", controllers.AddResourceActionToResourceRole)
	resourceActionRouter.DELETE("/:id", controllers.DeleteResourceRole)
	resourceActionRouter.PUT("/:id", controllers.UpdateResourceRole)
}
