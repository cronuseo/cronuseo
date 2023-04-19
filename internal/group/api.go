package group

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/:org_id/group")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
	// router.PATCH("/:id", res.patch)
}

type resource struct {
	service Service
}

// @Description Get group by ID.
// @Tags        Group
// @Param org_id path string true "Organization ID"
// @Param id path string true "Group ID"
// @Produce     json
// @Success     200 {object}  mongo_entity.Group
// @failure     404,500
// @Router      /{org_id}/group/{id} [get]
func (r resource) get(c echo.Context) error {

	group, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, group)
}

// @Description Get all groups.
// @Tags        Group
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  mongo_entity.Group
// @failure     500
// @Router      /{org_id}/group [get]
func (r resource) query(c echo.Context) error {

	var filter Filter
	org_id := c.Param("org_id")
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	groups, err := r.service.Query(c.Request().Context(), org_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusOK, groups)
}

// @Description Create group.
// @Tags        Group
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateGroupRequest true "body"
// @Produce     json
// @Success     201 {object}  mongo_entity.Group
// @failure     400,403,500
// @Router      /{org_id}/group [post]
func (r resource) create(c echo.Context) error {

	var input CreateGroupRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	group, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, group)
}

// @Description Update group.
// @Tags        Group
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Group ID"
// @Param request body UpdateGroupRequest true "body"
// @Produce     json
// @Success     201 {object}  mongo_entity.Group
// @failure     400,403,404,500
// @Router      /{org_id}/group/{id} [put]
func (r resource) update(c echo.Context) error {

	var input UpdateGroupRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	group, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, group)
}

// // @Description Delete group.
// // @Tags        Group
// // @Param org_id path string true "Organization ID"
// // @Param id path string true "Group ID"
// // @Produce     json
// // @Success     204
// // @failure     404,500
// // @Router      /{org_id}/group/{id} [delete]
func (r resource) delete(c echo.Context) error {

	err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
