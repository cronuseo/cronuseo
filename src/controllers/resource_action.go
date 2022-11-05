package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
	"net/http"
	"strconv"
)

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
	handlers.GetResourceActions(&resourceActions, resId)
	return c.JSON(http.StatusOK, &resourceActions)
}

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
	handlers.GetResourceAction(&resourceAction, resactId)
	return c.JSON(http.StatusOK, &resourceAction)

}

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
		if resourceAction.Key == "" || len(resourceAction.Key) < 4 || resourceAction.Name == "" || len(resourceAction.Name) < 4 {
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
	handlers.CreateResourceAction(&resourceAction)
	return c.JSON(http.StatusCreated, &resourceAction)
}

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
	handlers.DeleteResourceAction(&resourceAction, resactId)
	return c.JSON(http.StatusNoContent, "")
}

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
		if resourceAction.Key == "" || len(resourceAction.Key) < 4 || resourceAction.Name == "" || len(resourceAction.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&resourceAction); err != nil {
		return utils.InvalidErrorResponse()
	}
	handlers.UpdateResourceAction(&resourceAction, &reqResourceAction, resactId)
	return c.JSON(http.StatusCreated, &resourceAction)
}
