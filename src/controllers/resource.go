package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResources(c echo.Context) error {
	resources := []models.Resource{}
	proj_id := string(c.Param("proj_id"))
	exists, err := repositories.CheckProjectExistsById(proj_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	repositories.GetResources(&resources, proj_id)
	return c.JSON(http.StatusOK, &resources)
}

func GetResource(c echo.Context) error {
	var resource models.Resource
	res_id := string(c.Param("id"))
	proj_id := string(c.Param("proj_id"))
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})

	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	repositories.GetResource(&resource, res_id)
	return c.JSON(http.StatusOK, &resource)
}

func CreateResource(c echo.Context) error {
	var resource models.Resource
	proj_id := string(c.Param("proj_id"))
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	if err := c.Bind(&resource); err != nil {
		if resource.Key == "" || len(resource.Key) < 4 || resource.Name == "" || len(resource.Name) < 4 {
			return echo.NewHTTPError(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: err.Error()})
		}
	}
	if err := c.Validate(&resource); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: "Invalid inputs. Please check your inputs"})
	}
	int_proj_id, _ := strconv.Atoi(proj_id)
	resource.ProjectID = int_proj_id
	exists, err := repositories.CheckResourceExistsByKey(resource.Key, proj_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if exists {
		config.Log.Info("Resource already exists")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Resource already exists"})
	}
	repositories.CreateResource(&resource)
	return c.JSON(http.StatusOK, &resource)
}

func DeleteResource(c echo.Context) error {
	var resource models.Resource
	proj_id := string(c.Param("proj_id"))
	res_id := string(c.Param("id"))
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Resource not exists"})
	}
	repositories.DeleteResource(&resource, res_id)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateResource(c echo.Context) error {
	var resource models.Resource
	var reqResource models.Resource
	proj_id := string(c.Param("proj_id"))
	res_id := string(c.Param("id"))
	proj_exists, proj_err := repositories.CheckProjectExistsById(proj_id)
	if proj_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !proj_exists {
		config.Log.Info("Project not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Project not exists"})
	}
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Resource not exists"})
	}
	repositories.UpdateResource(&resource, &reqResource, res_id)
	return c.JSON(http.StatusCreated, &resource)
}
