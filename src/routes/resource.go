package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoutes(router *echo.Group) {

	resourceRouter := router.Group("/:proj_id/resources")

	resourceRouter.GET("", controllers.GetResources)
	resourceRouter.GET("/:id", controllers.GetResource)
	resourceRouter.POST("", controllers.CreateResource)
	resourceRouter.DELETE("/:id", controllers.DeleteResource)
	resourceRouter.PUT("/:id", controllers.UpdateResource)
}
