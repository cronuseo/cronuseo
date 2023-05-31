package check

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/cronuseo/internal/util"
)

type grpc_service struct {
	service Service
}

func (r grpc_service) grpcCheck(c echo.Context) error {
	var input CheckRequest
	api_key := c.Request().Header.Get("API_KEY")
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
	}

	allow, err := r.service.Check(context.Background(), c.Param("org"), input, api_key)
	if err != nil {
		return util.HandleError(err)
	}

	return c.JSON(http.StatusOK, allow)
}
