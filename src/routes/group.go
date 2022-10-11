package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func GroupRoutes(router *gin.Engine) {

	userRouter := router.Group("/groups")

	userRouter.GET("/:org_id", controllers.GetProjects)
	userRouter.GET("/:org_id/:id", controllers.GetProject)
	userRouter.POST("/:org_id", controllers.CreateProject)
	userRouter.DELETE("/:org_id/:id", controllers.Delete)
	userRouter.PUT("/:org_id/:id", controllers.UpdateProject)
}
