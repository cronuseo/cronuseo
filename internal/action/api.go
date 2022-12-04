package action

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := action{service}
	router := r.Group("/:resource_id/action")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type action struct {
	service Service
}

// @Description Get action by ID.
// @Tags        Action
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Action ID"
// @Produce     json
// @Success     200 {object}  entity.Action
// @failure     404,500
// @Router      /{resource_id}/action/{id} [get]
func (r action) get(c echo.Context) error {
	action, err := r.service.Get(c.Request().Context(), c.Param("resource_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, action)
}

// @Description Get all resources.
// @Tags        Action
// @Param resource_id path string true "Resource ID"
// @Produce     json
// @Success     200 {array}  entity.Action
// @failure     500
// @Router      /{resource_id}/action [get]
func (r action) query(c echo.Context) error {
	resources, err := r.service.Query(c.Request().Context(), c.Param("resource_id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusOK, resources)
}

// @Description Create action.
// @Tags        Action
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param request body CreateActionRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Action
// @failure     400,403,500
// @Router      /{resource_id}/action [post]
func (r action) create(c echo.Context) error {
	var input CreateActionRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	action, err := r.service.Create(c.Request().Context(), c.Param("resource_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, action)
}

// @Description Update action.
// @Tags        Action
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Action ID"
// @Param request body UpdateActionRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Action
// @failure     400,403,404,500
// @Router      /{resource_id}/action/{id} [put]
func (r action) update(c echo.Context) error {
	var input UpdateActionRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	action, err := r.service.Update(c.Request().Context(), c.Param("resource_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, action)
}

// @Description Delete action.
// @Tags        Action
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Action ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{resource_id}/action/{id} [delete]
func (r action) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("resource_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
