package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
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

func UserRoutes2(router *echo.Echo) {

	userRouter := router.Group("/users")

	userRouter.GET("/:org_id", controllers.GetUsers2)
	userRouter.GET("/:org_id/:id", controllers.GetUser2)
	userRouter.POST("/:org_id", controllers.CreateUser2)
	userRouter.DELETE("/:org_id/:id", controllers.DeleteUser2)
	userRouter.PUT("/:org_id/:id", controllers.UpdateUser2)
}
