package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/controllers"
)

func TenantRoutes(router *echo.Group) {

	tenantRouter := router.Group("/:org_id/tenant")

	tenantRouter.GET("", controllers.GetTenants)
	tenantRouter.GET("/:id", controllers.GetTenant)
	tenantRouter.POST("", controllers.CreateTenant)
	tenantRouter.DELETE("/:id", controllers.DeleteTenant)
	tenantRouter.PUT("/:id", controllers.UpdateTenant)
}
