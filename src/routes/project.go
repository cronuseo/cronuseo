package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ProjectRoutes(router *gin.Engine) {

	projectRouter := router.Group("/projects")

	projectRouter.GET("/:org_id", controllers.GetProjects)
	projectRouter.POST("/:org_id", controllers.CreateProjects)
	projectRouter.DELETE("/:org_id/:id", controllers.DeleteProjects)
	projectRouter.PUT("/:org_id/:id", controllers.UpdateProjects)
}
