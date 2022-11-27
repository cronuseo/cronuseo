package keto

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
)

func RegisterHandlers(r *echo.Group, service Service) {
	res := keto{service}
	router := r.Group("/:org/keto")
	router.POST("/create", res.create)
	router.POST("/check", res.check)
	router.POST("/list/object", res.getobjectlist)
	router.POST("/list/subject", res.getsubjectlist)
}

type keto struct {
	service Service
}

// @Description Create tuple.
// @Tags        Keto
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /org/keto/create [post]
func (r keto) create(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	err := r.service.CreateTuple(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, "")
}

// @Description Create tuple.
// @Tags        Keto
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /org/keto/check [post]
func (r keto) check(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckTuple(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Create tuple.
// @Tags        Keto
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /org/keto/delete [post]
func (r keto) delete(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckTuple(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Get objects.
// @Tags        Keto
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /org/keto/list/object [post]
func (r keto) getobjectlist(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.GetObjectListBySubject(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Get objects.
// @Tags        Keto
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /org/keto/list/subject [post]
func (r keto) getsubjectlist(c echo.Context) error {
	var input entity.Tuple
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.GetSubjectListByObject(context.Background(), c.Param("org"), "permission", input)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}
