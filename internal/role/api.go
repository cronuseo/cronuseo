package role

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := role{service}
	router := r.Group("/:org_id/role")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type role struct {
	service Service
}

// @Description Get role by ID.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Param id path string true "Role ID"
// @Produce     json
// @Success     200 {object}  entity.Role
// @failure     404,500
// @Router      /{org_id}/role/{id} [get]
func (r role) get(c echo.Context) error {
	role, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, role)
}

// @Description Get all roles.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  entity.Role
// @failure     500
// @Router      /{org_id}/role [get]
func (r role) query(c echo.Context) error {
	roles, err := r.service.Query(c.Request().Context(), c.Param("org_id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusOK, roles)
}

// @Description Create role.
// @Tags        Role
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateRoleRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Role
// @failure     400,403,500
// @Router      /{org_id}/role [post]
func (r role) create(c echo.Context) error {
	var input CreateRoleRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	role, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, role)
}

// @Description Update role.
// @Tags        Role
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Role ID"
// @Param request body UpdateRoleRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Role
// @failure     400,403,404,500
// @Router      /{org_id}/role/{id} [put]
func (r role) update(c echo.Context) error {
	var input UpdateRoleRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	role, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, role)
}

// @Description Delete role.
// @Tags        Role
// @Param org_id path string true "Organization ID"
// @Param id path string true "Role ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/role/{id} [delete]
func (r role) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
