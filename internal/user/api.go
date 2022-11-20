package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group(":org_id/user")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type resource struct {
	service Service
}

// @Description Get user by ID.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Produce     json
// @Success     200 {object}  entity.User
// @failure     404,500
// @Router      /{org_id}/user/{id} [get]
func (r resource) get(c echo.Context) error {
	user, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// @Description Get all users.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  entity.User
// @failure     500
// @Router      /{org_id}/user [get]
func (r resource) query(c echo.Context) error {
	users, err := r.service.Query(c.Request().Context(), c.Param("org_id"))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// @Description Create user.
// @Tags        User
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreateUserRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.User
// @failure     400,403,500
// @Router      /{org_id}/user [post]
func (r resource) create(c echo.Context) error {
	var input CreateUserRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	user, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		log.Debug(err.Error())
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

// @Description Update user.
// @Tags        User
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.User
// @failure     400,403,404,500
// @Router      /{org_id}/user/{id} [put]
func (r resource) update(c echo.Context) error {
	var input UpdateUserRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	user, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusCreated, user)
}

// @Description Delete user.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/user/{id} [delete]
func (r resource) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusNoContent, "")
}
