package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func UserRoutes(router *gin.Engine) {

	router.GET("/", controllers.GetUsers)
	router.POST("/", controllers.CreateUser)
	router.DELETE("/:id", controllers.DeleteUser)
	router.POST("/:id", controllers.UpdateUser)
}
