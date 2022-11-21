package resource

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/:org_id/resource")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type resource struct {
	service Service
}

// @Description Get resource by ID.
// @Tags        Resource
// @Param org_id path string true "Organization ID"
// @Param id path string true "Resource ID"
// @Produce     json
// @Success     200 {object}  entity.Resource
// @failure     404,500
// @Router      /{org_id}/resource/{id} [get]
func (r resource) get(c echo.Context) error {
	resource, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resource)
}

// @Description Get all resources.
// @Tags        Resource
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  entity.Resource
// @failure     500
// @Router      /{org_id}/resource [get]
func (r resource) query(c echo.Context) error {
	resources, err := r.service.Query(c.Request().Context(), c.Param("org_id"))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusOK, resources)
}

// @Description Create resource.
// @Tags        Resource
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateResourceRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Resource
// @failure     400,403,500
// @Router      /{org_id}/resource [post]
func (r resource) create(c echo.Context) error {
	var input CreateResourceRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	resource, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		log.Debug(err.Error())
		return err
	}

	return c.JSON(http.StatusCreated, resource)
}

// @Description Update resource.
// @Tags        Resource
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Resource ID"
// @Param request body UpdateResourceRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Resource
// @failure     400,403,404,500
// @Router      /{org_id}/resource/{id} [put]
func (r resource) update(c echo.Context) error {
	var input UpdateResourceRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	resource, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusCreated, resource)
}

// @Description Delete resource.
// @Tags        Resource
// @Param org_id path string true "Organization ID"
// @Param id path string true "Resource ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/resource/{id} [delete]
func (r resource) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusNoContent, "")
}
