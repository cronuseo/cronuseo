package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorResponse(code int, message string) *echo.HTTPError {
	return echo.NewHTTPError(code, message)
}

func ServerErrorResponse() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, "Server Error!")
}

func NotFoundErrorResponse(entity string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, entity+" not exists")
}

func InvalidErrorResponse() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs")
}

func AlreadyExistsErrorResponse(entity string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusForbidden, entity+" already exists")
}
