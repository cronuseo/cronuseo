package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func GroupRoutes(router *echo.Group) {

	userRouter := router.Group("/:org_id/group")

	userRouter.GET("", controllers.GetGroups)
	userRouter.GET("/:id", controllers.GetGroup)
	userRouter.POST("", controllers.CreateGroup)
	userRouter.POST("/:id/user", controllers.AddUsersToGroup)
	userRouter.DELETE("/:id", controllers.DeleteGroup)
	userRouter.PUT("/:id", controllers.UpdateGroup)
}
