package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func CheckRoutes(router *echo.Group) {

	checkRouter := router.Group("/check")

	checkRouter.POST("", controllers.CheckAllowed)
	checkRouter.POST("/list", controllers.Check)
}
