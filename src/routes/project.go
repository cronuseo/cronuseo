package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func ProjectRoutes(router *echo.Group) {

	projectRouter := router.Group("/:tenant_id/project")

	projectRouter.GET("", controllers.GetProjects)
	projectRouter.GET("/:id", controllers.GetProject)
	projectRouter.POST("", controllers.CreateProject)
	projectRouter.DELETE("/:id", controllers.DeleteProject)
	projectRouter.PUT("/:id", controllers.UpdateProject)
}
