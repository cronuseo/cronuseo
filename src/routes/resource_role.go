package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ResourceRoleRoutes(router *gin.Engine) {

	resourceActionRouter := router.Group("/resource_roles")

	resourceActionRouter.GET("/:res_id/:id", controllers.GetResourceRole)
	resourceActionRouter.GET("/:res_id", controllers.GetResourceRoles)
	resourceActionRouter.POST("/:res_id", controllers.CreateResourceRole)
	resourceActionRouter.POST("/:res_id/user/:user_id", controllers.CreateResourceRole)
	resourceActionRouter.POST("/:res_id/group/:group_id", controllers.CreateResourceRole)
	resourceActionRouter.DELETE("/:res_id/:id", controllers.DeleteResourceRole)
	resourceActionRouter.PUT("/:res_id/:id", controllers.UpdateResourceRole)
}
