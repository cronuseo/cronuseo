package test

import (
	"github.com/labstack/echo/v4"
)

// MockRouter creates a routing.Router for testing APIs.
func MockRouter() *echo.Echo {
	router := echo.New()
	return router
}
