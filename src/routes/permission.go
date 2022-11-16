package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func PermissionRoutes(router *echo.Group) {

	router.POST("/:orgKey/permission/:res_id", controllers.CreatePermission)
}
