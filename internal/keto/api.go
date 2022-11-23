package keto

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := keto{service}
	router := r.Group("/:org/keto")
	router.POST("", res.create)
}

type keto struct {
	service Service
}

// @Description Create tuple.
// @Tags        Keto
// @Accept      json
// @Param org path string true "Organization ID"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /org_id/keto [post]
func (r keto) create(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	err := r.service.CreateTuple(c.Request().Context(), "permission", input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "")
}
