package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceActionRoutes(router *gin.Engine) {

	resourceActionRouter := router.Group("/resource_actions")
	resourceActionRouter.GET("/:res_id", controllers.GetResourceActions)
	resourceActionRouter.GET("/:res_id/:id", controllers.GetResourceAction)
	resourceActionRouter.POST("/:res_id", controllers.CreateResourceAction)
	resourceActionRouter.DELETE("/:res_id/:id", controllers.DeleteResourceAction)
	resourceActionRouter.PUT("/:res_id/:id", controllers.UpdateResourceAction)
}

func ResourceActionRoutes2(router *echo.Echo) {

	resourceActionRouter := router.Group("/resource_actions")
	resourceActionRouter.GET("/:res_id", controllers.GetResourceActions2)
	resourceActionRouter.GET("/:res_id/:id", controllers.GetResourceAction2)
	resourceActionRouter.POST("/:res_id", controllers.CreateResourceAction2)
	resourceActionRouter.DELETE("/:res_id/:id", controllers.DeleteResourceAction2)
	resourceActionRouter.PUT("/:res_id/:id", controllers.UpdateResourceAction2)
}
