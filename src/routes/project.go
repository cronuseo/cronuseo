package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ProjectRoutes(router *gin.Engine) {

	router.GET("/:org_id/projects", controllers.GetProjects)
	router.POST("/:org_id/projects", controllers.CreateProjects)
	router.DELETE("/:org_id/projects/:id", controllers.DeleteProjects)
	router.PUT("/:org_id/projects/:id", controllers.UpdateProjects)
}
