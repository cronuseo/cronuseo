package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all resources.
// @Tags        Resource
// @Param proj_id path string true "Project ID"
// @Produce     json
// @Success     200 {array}  models.Resource
// @failure     500
// @Router      /{proj_id}/resource [get]
func GetResources(c echo.Context) error {
	resources := []models.Resource{}
	projId := string(c.Param("proj_id"))

	exists, _ := handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}

	err := handlers.GetResources(projId, &resources)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resources)
}

// @Description Get resource by ID.
// @Tags        Resource
// @Param proj_id path string true "Project ID"
// @Param id path string true "Resource ID"
// @Produce     json
// @Success     200 {object}  models.Resource
// @failure     404,500
// @Router      /{proj_id}/resource/{id} [get]
func GetResource(c echo.Context) error {
	var resource models.Resource
	resId := string(c.Param("id"))
	projId := string(c.Param("proj_id"))

	exists, _ := handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}

	exists, _ = handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	err := handlers.GetResource(projId, resId, &resource)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resource)
}

// @Description Create resource.
// @Tags        Resource
// @Accept      json
// @Param proj_id path string true "Project ID"
// @Param request body models.ResourceCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Resource
// @failure     400,403,500
// @Router      /{proj_id}/resource [post]
func CreateResource(c echo.Context) error {
	var resource models.Resource
	projId := string(c.Param("proj_id"))

	exists, _ := handlers.CheckProjectExistsById(projId)
	if !exists {
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

	resource.ProjectID = projId
	exists, _ = handlers.CheckResourceExistsByKey(resource.Key, projId)
	if exists {
		config.Log.Info("Resource already exists")
		return utils.AlreadyExistsErrorResponse("Resource")
	}

	err := handlers.CreateResource(projId, &resource)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resource)
}

// @Description Delete resource.
// @Tags        Resource
// @Param proj_id path string true "Project ID"
// @Param id path string true "Resource ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{proj_id}/resource/{id} [delete]
func DeleteResource(c echo.Context) error {

	projId := string(c.Param("proj_id"))
	resId := string(c.Param("id"))

	exists, _ := handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}

	exists, _ = handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	err := handlers.DeleteResource(projId, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update resource.
// @Tags        Resource
// @Accept      json
// @Param proj_id path string true "Project ID"
// @Param id path string true "Resource ID"
// @Param request body models.ResourceUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Resource
// @failure     400,403,404,500
// @Router      /{proj_id}/resource/{id} [put]
func UpdateResource(c echo.Context) error {
	var resource models.Resource
	var reqResource models.ResourceUpdateRequest
	projId := string(c.Param("proj_id"))
	resId := string(c.Param("id"))

	exists, _ := handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}

	exists, _ = handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	if err := c.Bind(&reqResource); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&reqResource); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.UpdateResource(projId, resId, &resource, &reqResource)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &resource)
}
