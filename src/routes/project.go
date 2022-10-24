package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ProjectRoutes(router *gin.Engine) {

	projectRouter := router.Group("/projects")

	projectRouter.GET("/:org_id", controllers.GetProjects)
	projectRouter.GET("/:org_id/:id", controllers.GetProject)
	projectRouter.POST("/:org_id", controllers.CreateProject)
	projectRouter.DELETE("/:org_id/:id", controllers.DeleteProject)
	projectRouter.PUT("/:org_id/:id", controllers.UpdateProject)
}

func ProjectRoutes2(router *echo.Echo) {

	projectRouter := router.Group("/projects")

	projectRouter.GET("/:org_id", controllers.GetProjects2)
	projectRouter.GET("/:org_id/:id", controllers.GetProject2)
	projectRouter.POST("/:org_id", controllers.CreateProject2)
	projectRouter.DELETE("/:org_id/:id", controllers.DeleteProject2)
	projectRouter.PUT("/:org_id/:id", controllers.UpdateProject2)
}
