package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func CheckRoutes(router *echo.Echo) {

	userRouter := router.Group("/check")

	userRouter.POST("", controllers.Check)
}
