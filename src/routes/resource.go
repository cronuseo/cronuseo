package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoutes(router *gin.Engine) {

	resourceRouter := router.Group("/resources")

	resourceRouter.GET("/:proj_id", controllers.GetResources)
	resourceRouter.GET("/:proj_id/:id", controllers.GetResource)
	resourceRouter.POST("/:proj_id", controllers.CreateResource)
	resourceRouter.DELETE("/:proj_id/:id", controllers.DeleteResource)
	resourceRouter.PUT("/:proj_id/:id", controllers.UpdateResource)
}

func ResourceRoutes2(router *echo.Echo) {

	resourceRouter := router.Group("/resources")

	resourceRouter.GET("/:proj_id", controllers.GetResources2)
	resourceRouter.GET("/:proj_id/:id", controllers.GetResource2)
	resourceRouter.POST("/:proj_id", controllers.CreateResource2)
	resourceRouter.DELETE("/:proj_id/:id", controllers.DeleteResource2)
	resourceRouter.PUT("/:proj_id/:id", controllers.UpdateResource2)
}
