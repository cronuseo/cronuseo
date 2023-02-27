package check

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/entity"
	"github.com/shashimalcse/cronuseo/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHandlers(r *echo.Group, service Service, monitoring_client *mongo.Collection) {
	res := permission{service: service, monitoring_client: monitoring_client}
	router := r.Group("/:org/permission/check")
	router.POST("", res.check)
	router.POST("/username", res.checkbyusername)
	router.POST("/multi_actions", res.checkpermissions)
	router.POST("/multi_resources", res.checkall)
}

type permission struct {
	service           Service
	monitoring_client *mongo.Collection
	wg                sync.WaitGroup
}

// @Description Check tuple.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body Tuple true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check [post]
func (r permission) check(c echo.Context) error {
	var input entity.Tuple
	api_key := c.Request().Header.Get("API_KEY")
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	allow, err := r.service.CheckTuple(context.Background(), c.Param("org"), "permission", input, api_key)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body entity.CheckRequestWithUser true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check/username [post]
func (r permission) checkbyusername(c echo.Context) error {
	var input entity.CheckRequestWithUser
	api_key := c.Request().Header.Get("API_KEY")
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	allow, err := r.service.CheckByUsername(context.Background(), c.Param("org"), "permission", input, api_key)
	if err != nil {
		return util.HandleError(err)
	}

	r.wg.Add(1)
	go func(input entity.CheckRequestWithUser) {
		defer r.wg.Done()
		checkMetrics := entity.CheckMetrics{Request: input, Result: allow, Timestamp: time.Now()}
		_, _ = r.monitoring_client.InsertOne(context.Background(), checkMetrics)
	}(input)

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body entity.CheckRequestWithPermissions true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check/multi_actions [post]
func (r permission) checkpermissions(c echo.Context) error {
	var input entity.CheckRequestWithPermissions
	api_key := c.Request().Header.Get("API_KEY")
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckPermissions(context.Background(), c.Param("org"), "permission", input, api_key)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}

// @Description Check by username.
// @Tags        Permission
// @Accept      json
// @Param org path string true "Organization"
// @Param request body entity.CheckRequestAll true "body"
// @Produce     json
// @Success     201
// @failure     400,403,500
// @Router      /{org}/permission/check/multi_resources [post]
func (r permission) checkall(c echo.Context) error {
	var input entity.CheckRequestAll
	api_key := c.Request().Header.Get("API_KEY")
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}
	allow, err := r.service.CheckAll(context.Background(), c.Param("org"), "permission", input, api_key)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}
