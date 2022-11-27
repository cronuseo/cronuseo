package permission

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := permission{service}
	router := r.Group("/:resource_id/permission")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type permission struct {
	service Service
}

// @Description Get permission by ID.
// @Tags        Permission
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Permission ID"
// @Produce     json
// @Success     200 {object}  entity.Permission
// @failure     404,500
// @Router      /{resource_id}/permission/{id} [get]
func (r permission) get(c echo.Context) error {
	permission, err := r.service.Get(c.Request().Context(), c.Param("resource_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, permission)
}

// @Description Get all resources.
// @Tags        Permission
// @Param resource_id path string true "Resource ID"
// @Produce     json
// @Success     200 {array}  entity.Permission
// @failure     500
// @Router      /{resource_id}/permission [get]
func (r permission) query(c echo.Context) error {
	resources, err := r.service.Query(c.Request().Context(), c.Param("resource_id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusOK, resources)
}

// @Description Create permission.
// @Tags        Permission
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param request body CreateResourceRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Permission
// @failure     400,403,500
// @Router      /{resource_id}/permission [post]
func (r permission) create(c echo.Context) error {
	var input CreateResourceRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	permission, err := r.service.Create(c.Request().Context(), c.Param("resource_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, permission)
}

// @Description Update permission.
// @Tags        Permission
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Permission ID"
// @Param request body UpdateResourceRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Permission
// @failure     400,403,404,500
// @Router      /{resource_id}/permission/{id} [put]
func (r permission) update(c echo.Context) error {
	var input UpdateResourceRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	permission, err := r.service.Update(c.Request().Context(), c.Param("resource_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, permission)
}

// @Description Delete permission.
// @Tags        Permission
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Permission ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{resource_id}/permission/{id} [delete]
func (r permission) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("resource_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
