package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func GroupRoutes(router *gin.Engine) {

	userRouter := router.Group("/groups")

	userRouter.GET("/:org_id", controllers.GetGroups)
	userRouter.GET("/:org_id/:id", controllers.GetGroup)
	userRouter.POST("/:org_id", controllers.CreateGroup)
	userRouter.POST("/:org_id/:id/:user_id", controllers.AddUserToGroup)
	userRouter.DELETE("/:org_id/:id", controllers.DeleteGroup)
	userRouter.PUT("/:org_id/:id", controllers.UpdateGroup)
}

func GroupRoutes2(router *echo.Echo) {

	userRouter := router.Group("/groups")

	userRouter.GET("/:org_id", controllers.GetGroups2)
	userRouter.GET("/:org_id/:id", controllers.GetGroup2)
	userRouter.POST("/:org_id", controllers.CreateGroup2)
	userRouter.POST("/:org_id/:id/:user_id", controllers.AddUserToGroup2)
	userRouter.DELETE("/:org_id/:id", controllers.DeleteGroup2)
	userRouter.PUT("/:org_id/:id", controllers.UpdateGroup2)
}
