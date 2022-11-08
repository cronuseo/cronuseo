package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all resource actions.
// @Tags        Resource Action
// @Param res_id path int true "Resource ID"
// @Produce     json
// @Success     200 {array}  models.Resource
// @failure     500
// @Router      /{res_id}/resource_action [get]
func GetResourceActions(c echo.Context) error {
	resourceActions := []models.ResourceAction{}
	resId := string(c.Param("res_id"))
	exists, err := handlers.CheckResourceExistsById(resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	err = handlers.GetResourceActions(&resourceActions, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceActions)
}

// @Description Get resource action by ID.
// @Tags        Resource Action
// @Param res_id path int true "Resource ID"
// @Param id path int true "Resource Action ID"
// @Produce     json
// @Success     200 {object}  models.ResourceAction
// @failure     404,500
// @Router      /{res_id}/resource_action/{id} [get]
func GetResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	resId := string(c.Param("res_id"))
	resactId := string(c.Param("id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resact_exists, resact_err := handlers.CheckResourceActionExistsById(resactId)
	if resact_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resact_exists {
		config.Log.Info("Resource Action not exists")
		return utils.NotFoundErrorResponse("Resource Action")
	}
	err := handlers.GetResourceAction(&resourceAction, resactId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceAction)

}

// @Description Create resource action.
// @Tags        Resource Action
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param request body models.ResourceActionCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceAction
// @failure     400,403,500
// @Router      /{res_id}/resource_action [post]
func CreateResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	resId := string(c.Param("res_id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	if err := c.Bind(&resourceAction); err != nil {
		if resourceAction.Key == "" || len(resourceAction.Key) < 4 ||
			resourceAction.Name == "" || len(resourceAction.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&resourceAction); err != nil {
		return utils.InvalidErrorResponse()
	}
	intResId, _ := strconv.Atoi(resId)
	resourceAction.ResourceID = intResId
	exists, err := handlers.CheckResourceActionExistsByKey(resourceAction.Key, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Resource Action already exists")
		return utils.AlreadyExistsErrorResponse("Resource Action")
	}
	err = handlers.CreateResourceAction(&resourceAction)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &resourceAction)
}

// @Description Delete resource action.
// @Tags        Resource Action
// @Param res_id path int true "Resource ID"
// @Param id path int true "Resource Action ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{res_id}/resource_action/{id} [delete]
func DeleteResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	resId := string(c.Param("res_id"))
	resactId := string(c.Param("id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resactExists, resactErr := handlers.CheckResourceActionExistsById(resactId)
	if resactErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resactExists {
		config.Log.Info("Resource Action not exists")
		return utils.NotFoundErrorResponse("Resource Action")
	}
	err := handlers.DeleteResourceAction(&resourceAction, resactId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update resource action.
// @Tags        Resource Action
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param id path int true "Resource Action ID"
// @Param request body models.ResourceActionUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceAction
// @failure     400,403,500
// @Router      /{proj_id}/resource_action/{id} [put]
func UpdateResourceAction(c echo.Context) error {
	var resourceAction models.ResourceAction
	var reqResourceAction models.ResourceAction
	resId := string(c.Param("res_id"))
	resactId := string(c.Param("id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resactExists, resactErr := handlers.CheckResourceActionExistsById(resactId)
	if resactErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resactExists {
		config.Log.Info("Resource Action not exists")
		return utils.NotFoundErrorResponse("Resource Action")
	}
	if err := c.Bind(&resourceAction); err != nil {
		if resourceAction.Key == "" || len(resourceAction.Key) < 4 ||
			resourceAction.Name == "" || len(resourceAction.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&resourceAction); err != nil {
		return utils.InvalidErrorResponse()
	}
	err := handlers.UpdateResourceAction(&resourceAction, &reqResourceAction, resactId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &resourceAction)
}
