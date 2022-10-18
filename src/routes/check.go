package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func CheckRoutes(router *gin.Engine) {

	userRouter := router.Group("/check")

	userRouter.POST("/", controllers.Check)
}
