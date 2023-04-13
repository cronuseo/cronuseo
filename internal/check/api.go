package check

import (
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Group, service Service) {
	// res := permission{service: service}
	// router := r.Group("/:org/permission/check")
	// router.POST("", res.check)
}

type permission struct {
	service Service
}
