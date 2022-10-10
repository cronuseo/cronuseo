package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoleRoutes(router *gin.Engine) {

	resourceActionRouter := router.Group("/resource_role")

	resourceActionRouter.GET("/:res_id", controllers.GetResource)
	resourceActionRouter.POST("/:res_id", controllers.CreateResource)
	resourceActionRouter.DELETE("/:res_id/:id", controllers.DeleteResource)
	resourceActionRouter.PUT("/:res_id/:id", controllers.UpdateResource)
}
