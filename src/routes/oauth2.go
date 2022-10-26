package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func OAuth2Routes(router *echo.Echo) {

	router.GET("/authorize", controllers.HandleAuthentication)
	router.GET("/callback", controllers.HandleCallback)
}
