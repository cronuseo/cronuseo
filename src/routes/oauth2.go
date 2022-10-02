package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func OAuth2Routes(router *gin.Engine) {

	router.GET("/authorize", controllers.HandleAuthentication)
	router.GET("/callback", controllers.HandleCallback)
}
