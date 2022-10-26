package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func UserRoutes(router *echo.Echo) {

	userRouter := router.Group("/users")

	userRouter.GET("/:org_id", controllers.GetUsers)
	userRouter.GET("/:org_id/:id", controllers.GetUser)
	userRouter.POST("/:org_id", controllers.CreateUser)
	userRouter.DELETE("/:org_id/:id", controllers.DeleteUser)
	userRouter.PUT("/:org_id/:id", controllers.UpdateUser)
}
