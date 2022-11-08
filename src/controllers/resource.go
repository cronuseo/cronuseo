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

// @Description Get all resources.
// @Tags        Resource
// @Param proj_id path int true "Project ID"
// @Produce     json
// @Success     200 {array}  models.Resource
// @failure     500
// @Router      /{proj_id}/resource [get]
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
	err = handlers.GetResources(&resources, projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resources)
}

// @Description Get resource by ID.
// @Tags        Resource
// @Param proj_id path int true "Project ID"
// @Param id path int true "Resource ID"
// @Produce     json
// @Success     200 {object}  models.Resource
// @failure     404,500
// @Router      /{proj_id}/resource/{id} [get]
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
	err := handlers.GetResource(&resource, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resource)
}

// @Description Create resource.
// @Tags        Resource
// @Accept      json
// @Param proj_id path int true "Project ID"
// @Param request body models.ResourceCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Resource
// @failure     400,403,500
// @Router      /{proj_id}/resource [post]
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
	err = handlers.CreateResource(&resource)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resource)
}

// @Description Delete resource.
// @Tags        Resource
// @Param proj_id path int true "Project ID"
// @Param id path int true "Resource ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{proj_id}/resource/{id} [delete]
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
	err := handlers.DeleteResource(&resource, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update resource.
// @Tags        Resource
// @Accept      json
// @Param proj_id path int true "Project ID"
// @Param id path int true "Resource ID"
// @Param request body models.ResourceUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Resource
// @failure     400,403,500
// @Router      /{proj_id}/resource/{id} [put]
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
	err := handlers.UpdateResource(&resource, &reqResource, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &resource)
}
