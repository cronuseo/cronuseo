package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func UserRoutes(router *gin.Engine) {

	userRouter := router.Group("/users")

	userRouter.GET("/:org_id", controllers.GetUsers)
	userRouter.GET("/:org_id/:id", controllers.GetUser)
	userRouter.POST("/:org_id", controllers.CreateUser)
	userRouter.DELETE("/:org_id/:id", controllers.DeleteUser)
	userRouter.PUT("/:org_id/:id", controllers.UpdateUser)
}
