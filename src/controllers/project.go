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
	handlers.GetProjects(&projects, orgId)
	return c.JSON(http.StatusOK, &projects)
}

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
	handlers.GetProject(&proj, projId)
	return c.JSON(http.StatusOK, &proj)

}

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
	handlers.CreateProject(&project)
	return c.JSON(http.StatusCreated, &project)

}

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
	handlers.DeleteProject(&project, projId)
	return c.JSON(http.StatusNoContent, "")
}

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
	handlers.UpdateProject(&project, &reqProject, projId)
	return c.JSON(http.StatusCreated, &project)
}
