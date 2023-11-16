package policy

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}
	router := r.Group("/o/:org_id/policies")
	router.GET("", res.query)
	router.GET("/:id", res.get)
	router.POST("", res.create)
	router.DELETE("/:id", res.delete)
	router.PUT("/:id", res.update)
	router.PATCH("/:id", res.patch)
}

type resource struct {
	service Service
}

// @Description Get policy by ID.
// @Tags        Policy
// @Param org_id path string true "Organization ID"
// @Param id path string true "Policy ID"
// @Produce     json
// @Success     200 {object}  Policy
// @failure     404,500
// @Router      /{org_id}/polcies/{id} [get]
func (r resource) get(c echo.Context) error {

	policy, err := r.service.Get(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, policy)
}

// @Description Get all policies.
// @Tags        Policy
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  Policy
// @failure     500
// @Router      /{org_id}/polcies [get]
func (r resource) query(c echo.Context) error {

	var filter Filter
	org_id := c.Param("org_id")
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	users, err := r.service.Query(c.Request().Context(), org_id, filter)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusOK, users)
}

// @Description Create policy.
// @Tags        Policy
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body CreatePolicyRequest true "body"
// @Produce     json
// @Success     201 {object}  Policy
// @failure     400,403,500
// @Router      /{org_id}/polcies [post]
func (r resource) create(c echo.Context) error {

	var input CreatePolicyRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	user, err := r.service.Create(c.Request().Context(), c.Param("org_id"), input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusCreated, user)
}

// @Description Update policy.
// @Tags        Policy
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Policy ID"
// @Param request body UpdatePolicyRequest true "body"
// @Produce     json
// @Success     201 {object}  Policy
// @failure     400,403,404,500
// @Router      /{org_id}/polcies/{id} [put]
func (r resource) update(c echo.Context) error {

	var input UpdatePolicyRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	user, err := r.service.Update(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, user)
}

// @Description Patch policy.
// @Tags        Policy
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Policy ID"
// @Param request body PatchPolicyRequest true "body"
// @Produce     json
// @Success     201 {object}  Policy
// @failure     400,403,404,500
// @Router      /{org_id}/polcies/{id} [patch]
func (r resource) patch(c echo.Context) error {

	var input PatchPolicyRequest
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	user, err := r.service.Patch(c.Request().Context(), c.Param("org_id"), c.Param("id"), input)
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusCreated, user)
}

// // @Description Delete policy.
// // @Tags        Policy
// // @Param org_id path string true "Organization ID"
// // @Param id path string true "Policy ID"
// // @Produce     json
// // @Success     204
// // @failure     404,500
// // @Router      /{org_id}/policies/{id} [delete]
func (r resource) delete(c echo.Context) error {

	err := r.service.Delete(c.Request().Context(), c.Param("org_id"), c.Param("id"))
	if err != nil {
		return util.HandleError(err)
	}
	return c.JSON(http.StatusNoContent, "")
}
