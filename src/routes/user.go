package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func UserRoutes(router *gin.Engine) {

	router.GET("/users", controllers.GetUsers)
	router.POST("/users", controllers.CreateUser)
	router.DELETE("/users/:id", controllers.DeleteUser)
	router.POST("/users/:id", controllers.UpdateUser)
}
