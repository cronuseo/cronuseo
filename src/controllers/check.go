package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Check RBAC.
// @Tags        Check
// @Accept      json
// @Param request body models.ResourceRoleToResourceActionKey true "body"
// @Produce     json
// @Success     200 {string}  allowed    "status"
// @failure     403,500
// @Router      /check [post]
func CheckAllowed(c echo.Context) error {
	var keys models.ResourceRoleToResourceActionKey
	if err := c.Bind(&keys); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&keys); err != nil {
		return utils.InvalidErrorResponse()
	}
	allow, err := handlers.CheckAllowed(keys.Resource, keys.ResourceRole, keys.ResourceAction)
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

// @Description Check RBAC list.
// @Tags        Check
// @Accept      json
// @Param request body models.ResourceRoleToResourceActionKey true "body"
// @Produce     json
// @Success     200 {string}  read,write    "scopes"
// @failure     403,500
// @Router      /check/list [post]
func Check(c echo.Context) error {
	var keys models.ResourceRoleToResourceActionKey
	if err := c.Bind(&keys); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&keys); err != nil {
		return utils.InvalidErrorResponse()
	}
	allowedActions := handlers.Check(keys.Resource, keys.ResourceRole, keys.ResourceAction)

	return c.JSON(http.StatusOK, allowedActions)

}
