package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceActionRoutes(router *gin.Engine) {

	resourceActionRouter := router.Group("/resource_action")
	resourceActionRouter.GET("/:res_id", controllers.GetResources)
	resourceActionRouter.GET("/:res_id/:id", controllers.GetResource)
	resourceActionRouter.POST("/:res_id", controllers.CreateResource)
	resourceActionRouter.DELETE("/:res_id/:id", controllers.DeleteResource)
	resourceActionRouter.PUT("/:res_id/:id", controllers.UpdateResource)
}
