package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func CheckRoutes(router *gin.Engine) {

	userRouter := router.Group("/check")

	userRouter.POST("/", controllers.Check)
}

func CheckRoutes2(router *echo.Echo) {

	userRouter := router.Group("/check")

	userRouter.POST("/", controllers.Check2)
}
