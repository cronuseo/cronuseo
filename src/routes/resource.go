package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoutes(router *echo.Echo) {

	resourceRouter := router.Group("/resources")

	resourceRouter.GET("/:proj_id", controllers.GetResources)
	resourceRouter.GET("/:proj_id/:id", controllers.GetResource)
	resourceRouter.POST("/:proj_id", controllers.CreateResource)
	resourceRouter.DELETE("/:proj_id/:id", controllers.DeleteResource)
	resourceRouter.PUT("/:proj_id/:id", controllers.UpdateResource)
}
