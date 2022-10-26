package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetProjects(c *gin.Context) {
	projects := []models.Project{}
	org_id := string(c.Param("org_id"))
	exists, err := repositories.CheckOrganizationExistsById(org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	repositories.GetProjects(&projects, org_id)
	c.JSON(http.StatusOK, &projects)
}

func GetProjects2(c echo.Context) error {
	projects := []models.Project{}
	org_id := string(c.Param("org_id"))
	exists, err := repositories.CheckOrganizationExistsById(org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	repositories.GetProjects(&projects, org_id)
	return c.JSON(http.StatusOK, &projects)
}

func GetProject(c *gin.Context) {
	var proj models.Project
	org_id := string(c.Param("org_id"))
	proj_id := string(c.Param("id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
		return
	}
	repositories.GetProject(&proj, proj_id)
	c.JSON(http.StatusOK, &proj)

}

func GetProject2(c echo.Context) error {
	var proj models.Project
	org_id := string(c.Param("org_id"))
	proj_id := string(c.Param("id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	repositories.GetProject(&proj, proj_id)
	return c.JSON(http.StatusOK, &proj)

}

func CreateProject(c *gin.Context) {
	var project models.Project
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	if err := c.ShouldBindJSON(&project); err != nil {
		if project.Key == "" || len(project.Key) < 4 || project.Name == "" || len(project.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	int_org_id, _ := strconv.Atoi(org_id)
	project.OrganizationID = int_org_id
	exists, err := repositories.CheckProjectExistsByKey(project.Key, org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Project already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project already exists"})
		return
	}
	repositories.CreateProject(&project)
	c.JSON(http.StatusOK, &project)

}

func CreateProject2(c echo.Context) error {
	var project models.Project
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
	}
	if err := c.Bind(&project); err != nil {
		if project.Key == "" || len(project.Key) < 4 || project.Name == "" || len(project.Name) < 4 {
			return echo.NewHTTPError(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: err.Error()})
		}
	}
	if err := c.Validate(&project); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: "Invalid inputs. Please check your inputs"})
	}
	int_org_id, _ := strconv.Atoi(org_id)
	project.OrganizationID = int_org_id
	exists, err := repositories.CheckProjectExistsByKey(project.Key, org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if exists {
		config.Log.Info("Project already exists")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Project already exists"})
	}
	repositories.CreateProject(&project)
	return c.JSON(http.StatusCreated, &project)

}

func DeleteProject(c *gin.Context) {
	var project models.Project
	proj_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
		return
	}
	repositories.DeleteProject(&project, proj_id)
	c.JSON(http.StatusOK, "")
}

func DeleteProject2(c echo.Context) error {
	var project models.Project
	proj_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	repositories.DeleteProject(&project, proj_id)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateProject(c *gin.Context) {
	var project models.Project
	var reqProject models.Project
	proj_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
		return
	}
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
		return
	}
	repositories.UpdateProject(&project, &reqProject, proj_id)
	c.JSON(http.StatusOK, &project)
}

func UpdateProject2(c echo.Context) error {
	var project models.Project
	var reqProject models.Project
	proj_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
	}
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Project not exists"})
	}
	if err := c.Bind(&reqProject); err != nil {
		if reqProject.Name == "" {
			return echo.NewHTTPError(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: "Invalid inputs. Please check your inputs"})
		}
	}
	repositories.UpdateProject(&project, &reqProject, proj_id)
	return c.JSON(http.StatusCreated, &project)
}
