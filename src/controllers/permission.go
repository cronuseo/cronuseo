package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Create permission.
// @Tags        Permission
// @Accept      json
// @Param orgKey path string true "Organization Key"
// @Param res_id path string true "Resource ID"
// @Param request body models.Permissions true "body"
// @Produce     json
// @Success     200
// @failure     400,403,500
// @Router      /{orgKey}/permission/{res_id} [post]
func CreatePermission(c echo.Context) error {
	var permissions models.Permissions
	resourceId := string(c.Param("res_id"))
	orgKey := string(c.Param("orgKey"))

	exists, _ := handlers.CheckOrganizationExistsByKey(orgKey)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	if err := c.Bind(&permissions); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&permissions); err != nil {
		return utils.InvalidErrorResponse()
	}

	exists, _ = handlers.CheckResourceExistsById(resourceId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	err := handlers.CreatePermissions(orgKey, resourceId, &permissions)
	if err != nil {
		config.Log.Panic(err)
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, "")
}
