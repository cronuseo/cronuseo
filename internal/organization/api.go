package organization

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/organization")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
}

type resource struct {
	service Service
}

// @Description Get organization by ID.
// @Tags        Organization
// @Param id path string true "Organization ID"
// @Produce     json
// @Success     200 {object}  entity.Organization
// @failure     404,500
// @Router      /organization/{id} [get]
func (r resource) get(c echo.Context) error {
	organization, err := r.service.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, organization)
}

// @Description Get all organizations.
// @Tags        Organization
// @Produce     json
// @Success     200 {array}  entity.Organization
// @failure     500
// @Router      /organization [get]
func (r resource) query(c echo.Context) error {
	organizations, err := r.service.Query(c.Request().Context())
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusOK, organizations)
}

// @Description Create organization.
// @Tags        Organization
// @Accept      json
// @Param request body CreateOrganizationRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Organization
// @failure     400,403,500
// @Router      /organization [post]
func (r resource) create(c echo.Context) error {
	var input CreateOrganizationRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	organization, err := r.service.Create(c.Request().Context(), input)
	if err != nil {
		log.Debug(err.Error())
		return err
	}

	return c.JSON(http.StatusCreated, organization)
}

// @Description Update organization.
// @Tags        Organization
// @Accept      json
// @Param id path string true "Organization ID"
// @Param request body UpdateOrganizationRequest true "body"
// @Produce     json
// @Success     201 {object}  entity.Organization
// @failure     400,403,404,500
// @Router      /organization/{id} [put]
func (r resource) update(c echo.Context) error {
	var input UpdateOrganizationRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	organization, err := r.service.Update(c.Request().Context(), c.Param("id"), input)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusCreated, organization)
}

// @Description Delete organization.
// @Tags        Organization
// @Param id path string true "Organization ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /organization/{id} [delete]
func (r resource) delete(c echo.Context) error {
	_, err := r.service.Delete(c.Request().Context(), c.Param("id"))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return c.JSON(http.StatusNoContent, "")
}
