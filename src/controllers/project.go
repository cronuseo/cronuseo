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

// @Description Get all projects.
// @Tags        Project
// @Param org_id path int true "Organization ID"
// @Produce     json
// @Success     200 {array}  models.Project
// @failure     500
// @Router      /{org_id}/project [get]
func GetProjects(c echo.Context) error {
	projects := []models.Project{}
	orgId := string(c.Param("org_id"))
	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	err = handlers.GetProjects(&projects, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &projects)
}

// @Description Get project by ID.
// @Tags        Project
// @Param org_id path int true "Organization ID"
// @Param id path int true "Project ID"
// @Produce     json
// @Success     200 {object}  models.Project
// @failure     404,500
// @Router      /{org_id}/project/{id} [get]
func GetProject(c echo.Context) error {
	var proj models.Project
	orgId := string(c.Param("org_id"))
	projId := string(c.Param("id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	err := handlers.GetProject(&proj, projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &proj)

}

// @Description Create project.
// @Tags        Project
// @Accept      json
// @Param org_id path int true "Organization ID"
// @Param request body models.ProjectCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Project
// @failure     400,403,500
// @Router      /{org_id}/project [post]
func CreateProject(c echo.Context) error {
	var project models.Project
	orgId := string(c.Param("org_id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	if err := c.Bind(&project); err != nil {
		if project.Key == "" || len(project.Key) < 4 || project.Name == "" || len(project.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&project); err != nil {
		return utils.InvalidErrorResponse()
	}
	intOrgId, _ := strconv.Atoi(orgId)
	project.OrganizationID = intOrgId
	exists, err := handlers.CheckProjectExistsByKey(project.Key, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Project already exists")
		return utils.AlreadyExistsErrorResponse("Project")
	}
	err = handlers.CreateProject(&project)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &project)

}

// @Description Delete project.
// @Tags        Project
// @Param org_id path int true "Organization ID"
// @Param id path int true "Project ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/project/{id} [delete]
func DeleteProject(c echo.Context) error {
	var project models.Project
	projId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	err := handlers.DeleteProject(&project, projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update project.
// @Tags        Project
// @Accept      json
// @Param org_id path int true "Organization ID"
// @Param id path int true "Project ID"
// @Param request body models.ProjectUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Project
// @failure     400,403,404,500
// @Router      /{org_id}/project/{id} [put]
func UpdateProject(c echo.Context) error {
	var project models.Project
	var reqProject models.Project
	projId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	projExists, projErr := handlers.CheckProjectExistsById(projId)
	if projErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !projExists {
		config.Log.Info("Project not exists")
		return utils.NotFoundErrorResponse("Project")
	}
	if err := c.Bind(&reqProject); err != nil {
		if reqProject.Name == "" {
			return utils.ServerErrorResponse()
		}
	}
	err := handlers.UpdateProject(&project, &reqProject, projId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &project)
}
