package monitoring

import (
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/auth"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/:org_id/monitoring")
	config := echojwt.Config{
		SigningKey: []byte(auth.SecretKey),
	}
	router.Use(echojwt.WithConfig(config))
	router.GET("/allowed_data", res.getAllowedData)
}

type resource struct {
	service Service
}

// @Description Get allowed by org Id.
// @Tags        Organization
// @Param id path string true "Organization ID"
// @Produce     json
// @Success     200 {object}  entity.Organization
// @failure     404,500
// @Router      {org_id}/monitoring/allowed_data [get]
func (r resource) getAllowedData(c echo.Context) error {
	organization, err := r.service.GetAllowedData(c.Request().Context(), c.Param("org_id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, organization)
}
