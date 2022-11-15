package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all resource actions.
// @Tags        Resource Action
// @Param resource_id path string true "Resource ID"
// @Produce     json
// @Success     200 {array}  models.ResourceAction
// @failure     500
// @Router      /{res_id}/resource_action [get]
func GetResourceActions(c echo.Context) error {
	resourceActions := []models.ResourceAction{}
	resourceId := string(c.Param("res_id"))

	exists, _ := handlers.CheckResourceExistsById(resourceId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	err := handlers.GetResourceActions(resourceId, &resourceActions)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceActions)
}

// @Description Get resource action by ID.
// @Tags        Resource Action
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Resource Action ID"
// @Produce     json
// @Success     200 {object}  models.ResourceAction
// @failure     404,500
// @Router      /{res_id}/resource_action/{id} [get]
func GetResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	resourceId := string(c.Param("res_id"))
	resourceActionId := string(c.Param("id"))

	exists, _ := handlers.CheckResourceExistsById(resourceId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	exists, _ = handlers.CheckResourceActionExistsById(resourceActionId)
	if !exists {
		config.Log.Info("Resource Action not exists")
		return utils.NotFoundErrorResponse("Resource Action")
	}

	err := handlers.GetResourceAction(resourceId, resourceActionId, &resourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceAction)

}

// @Description Create resource action.
// @Tags        Resource Action
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param request body models.ResourceActionCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceAction
// @failure     400,403,500
// @Router      /{res_id}/resource_action [post]
func CreateResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	resourceId := string(c.Param("res_id"))

	exists, _ := handlers.CheckResourceExistsById(resourceId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	if err := c.Bind(&resourceAction); err != nil {
		if resourceAction.Name == "" || len(resourceAction.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&resourceAction); err != nil {
		return utils.InvalidErrorResponse()
	}

	resourceAction.ResourceID = resourceId
	exists, _ = handlers.CheckResourceActionExistsByKey(resourceId, resourceAction.Key)
	if exists {
		config.Log.Info("Resource Action already exists")
		return utils.AlreadyExistsErrorResponse("Resource Action")
	}

	err := handlers.CreateResourceAction(resourceId, &resourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &resourceAction)

}

// @Description Delete resource action.
// @Tags        Resource Action
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Resource Action ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{res_id}/resource_action/{id} [delete]
func DeleteResourceAction(c echo.Context) error {

	resourceActionId := string(c.Param("id"))
	resourceId := string(c.Param("res_id"))

	exists, _ := handlers.CheckResourceExistsById(resourceId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	exists, _ = handlers.CheckResourceActionExistsById(resourceActionId)
	if !exists {
		config.Log.Info("Resource Action not exists")
		return utils.NotFoundErrorResponse("Resource Action")
	}

	err := handlers.DeleteResourceAction(resourceId, resourceActionId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update resource action.
// @Tags        Resource Action
// @Accept      json
// @Param resource_id path string true "Resource ID"
// @Param id path string true "Resource Action ID"
// @Param request body models.ResourceActionUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceAction
// @failure     400,403,404,500
// @Router      /{res_id}/resource_action/{id} [put]
func UpdateResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	var reqResourceAction models.ResourceActionUpdateRequest
	resourceActionId := string(c.Param("id"))
	resourceId := string(c.Param("res_id"))

	exists, _ := handlers.CheckResourceExistsById(resourceId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	exists, _ = handlers.CheckResourceActionExistsById(resourceActionId)
	if !exists {
		config.Log.Info("ResourceAction not exists")
		return utils.NotFoundErrorResponse("ResourceAction")
	}

	if err := c.Bind(&reqResourceAction); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&reqResourceAction); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.UpdateResourceAction(resourceId, resourceActionId, &resourceAction, &reqResourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &resourceAction)
}
