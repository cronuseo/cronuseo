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

func GetResources(c echo.Context) error {
	resources := []models.Resource{}
	projId := string(c.Param("proj_id"))
	exists, err := handlers.CheckProjectExistsById(projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	handlers.GetResources(&resources, projId)
	return c.JSON(http.StatusOK, &resources)
}

func GetResource(c echo.Context) error {
	var resource models.Resource
	resId := string(c.Param("id"))
	projId := string(c.Param("proj_id"))
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()

	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	handlers.GetResource(&resource, resId)
	return c.JSON(http.StatusOK, &resource)
}

func CreateResource(c echo.Context) error {
	var resource models.Resource
	projId := string(c.Param("proj_id"))
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	if err := c.Bind(&resource); err != nil {
		if resource.Key == "" || len(resource.Key) < 4 || resource.Name == "" || len(resource.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&resource); err != nil {
		return utils.InvalidErrorResponse()
	}
	intProjId, _ := strconv.Atoi(projId)
	resource.ProjectID = intProjId
	exists, err := handlers.CheckResourceExistsByKey(resource.Key, projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Resource already exists")
		return utils.AlreadyExistsErrorResponse("Resource")
	}
	handlers.CreateResource(&resource)
	return c.JSON(http.StatusOK, &resource)
}

func DeleteResource(c echo.Context) error {
	var resource models.Resource
	projId := string(c.Param("proj_id"))
	resId := string(c.Param("id"))
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	handlers.DeleteResource(&resource, resId)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateResource(c echo.Context) error {
	var resource models.Resource
	var reqResource models.Resource
	projId := string(c.Param("proj_id"))
	resId := string(c.Param("id"))
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	handlers.UpdateResource(&resource, &reqResource, resId)
	return c.JSON(http.StatusCreated, &resource)
}
