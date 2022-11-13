package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all projects.
// @Tags        Project
// @Param tenant_id path string true "Tenant ID"
// @Produce     json
// @Success     200 {array}  models.Project
// @failure     500
// @Router      /{tenant_id}/project [get]
func GetProjects(c echo.Context) error {
	projects := []models.Project{}
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	err := handlers.GetProjects(tenantId, &projects)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &projects)
}

// @Description Get project by ID.
// @Tags        Project
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Project ID"
// @Produce     json
// @Success     200 {object}  models.Project
// @failure     404,500
// @Router      /{tenant_id}/project/{id} [get]
func GetProject(c echo.Context) error {
	var proj models.Project
	tenantId := string(c.Param("tenant_id"))
	projId := string(c.Param("id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	exists, _ = handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}

	err := handlers.GetProject(tenantId, projId, &proj)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &proj)

}

// @Description Create project.
// @Tags        Project
// @Accept      json
// @Param tenant_id path string true "Tenant ID"
// @Param request body models.ProjectCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Project
// @failure     400,403,500
// @Router      /{tenant_id}/project [post]
func CreateProject(c echo.Context) error {
	var project models.Project
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	if err := c.Bind(&project); err != nil {
		if project.Key == "" || len(project.Key) < 4 || project.Name == "" || len(project.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&project); err != nil {
		return utils.InvalidErrorResponse()
	}

	project.TenantID = tenantId
	exists, _ = handlers.CheckProjectExistsByKey(tenantId, project.Key)
	if exists {
		config.Log.Info("Project already exists")
		return utils.AlreadyExistsErrorResponse("Project")
	}

	err := handlers.CreateProject(tenantId, &project)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &project)

}

// @Description Delete project.
// @Tags        Project
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Project ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{tenant_id}/project/{id} [delete]
func DeleteProject(c echo.Context) error {

	projId := string(c.Param("id"))
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	exists, _ = handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}

	err := handlers.DeleteProject(tenantId, projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update project.
// @Tags        Project
// @Accept      json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Project ID"
// @Param request body models.ProjectUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Project
// @failure     400,403,404,500
// @Router      /{tenant_id}/project/{id} [put]
func UpdateProject(c echo.Context) error {
	var project models.Project
	var reqProject models.Project
	projId := string(c.Param("id"))
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	exists, _ = handlers.CheckProjectExistsById(projId)
	if !exists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	if err := c.Bind(&reqProject); err != nil {
		if reqProject.Name == "" {
			return utils.ServerErrorResponse()
		}
	}

	err := handlers.UpdateProject(tenantId, projId, &project, &reqProject)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &project)
}
