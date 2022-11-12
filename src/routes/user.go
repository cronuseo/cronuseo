package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func UserRoutes(router *echo.Group) {

	userRouter := router.Group("/:tenant_id/user")

	userRouter.GET("", controllers.GetUsers)
	userRouter.GET("/:id", controllers.GetUser)
	userRouter.POST("", controllers.CreateUser)
	userRouter.DELETE("/:id", controllers.DeleteUser)
	userRouter.PUT("/:id", controllers.UpdateUser)
}
