package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
	"net/http"
)

func Check(c echo.Context) error {
	var keys models.ResourceRoleToResourceActionKey
	if err := c.Bind(&keys); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&keys); err != nil {
		return utils.InvalidErrorResponse()
	}
	allow, err := handlers.Check(keys.Resource, keys.ResourceRole, keys.ResourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if allow {
		return c.JSON(http.StatusOK, "allowed")
	} else {
		return c.JSON(http.StatusOK, "not allowed")
	}

}
