package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func GroupRoutes(router *echo.Group) {

	userRouter := router.Group("/:org_id/groups")

	userRouter.GET("", controllers.GetGroups)
	userRouter.GET("/:id", controllers.GetGroup)
	userRouter.POST("", controllers.CreateGroup)
	userRouter.POST("/:id/:user_id", controllers.AddUserToGroup)
	userRouter.DELETE("/:id", controllers.DeleteGroup)
	userRouter.PUT("/:id", controllers.UpdateGroup)
}
