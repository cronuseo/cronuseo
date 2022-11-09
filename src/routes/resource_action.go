package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceActionRoutes(router *echo.Group) {

	resourceActionRouter := router.Group("/:res_id/resource_actions")
	resourceActionRouter.GET("", controllers.GetResourceActions)
	resourceActionRouter.GET("/:id", controllers.GetResourceAction)
	resourceActionRouter.POST("", controllers.CreateResourceAction)
	resourceActionRouter.DELETE("/:id", controllers.DeleteResourceAction)
	resourceActionRouter.PUT("/:id", controllers.UpdateResourceAction)
}
