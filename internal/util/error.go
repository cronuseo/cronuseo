package util

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AlreadyExistsError struct {
	Path string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%v already exists.", e.Path)
}

type NotFoundError struct {
	Path string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v not found.", e.Path)
}

type InvalidInputError struct {
	Path string
}

type SystemError struct {
	Message string
}

func (e *SystemError) Error() string {
	return e.Message
}

func (e *InvalidInputError) Error() string {
	if e.Path != "" {
		return "Invalid " + e.Path
	}
	return "Invalid input."
}

func HandleError(err error) *echo.HTTPError {
	switch e := err.(type) {
	case *InvalidInputError:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid inputs. Please check your inputs.")
	case *AlreadyExistsError:
		return echo.NewHTTPError(http.StatusConflict, e.Error())
	case *NotFoundError:
		return echo.NewHTTPError(http.StatusNotFound, e.Error())
	case *SystemError:
		return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Server Error!")
	}
}
