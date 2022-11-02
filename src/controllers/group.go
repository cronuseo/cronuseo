package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetGroups(c echo.Context) error {
	groups := []models.Group{}
	org_id := string(c.Param("org_id"))
	exists, err := repositories.CheckOrganizationExistsById(org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	repositories.GetGroups(&groups, org_id)
	return c.JSON(http.StatusOK, &groups)
}

func GetGroup(c echo.Context) error {
	var group models.GroupUsers
	org_id := string(c.Param("org_id"))
	group_id := string(c.Param("id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	group_exists, group_err := repositories.CheckGroupExistsById(group_id)
	if group_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !group_exists {
		config.Log.Info("Group not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Group not exists"})
	}
	repositories.GetUsersFromGroup(group_id, &group)
	return c.JSON(http.StatusOK, &group)

}

func CreateGroup(c echo.Context) error {
	var group models.Group
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
	}
	if err := c.Bind(&group); err != nil {
		if group.Key == "" || len(group.Key) < 4 || group.Name == "" || len(group.Name) < 4 {
			return echo.NewHTTPError(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: err.Error()})
		}
	}
	if err := c.Validate(&group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: "Invalid inputs. Please check your inputs"})
	}
	int_org_id, _ := strconv.Atoi(org_id)
	group.OrganizationID = int_org_id
	exists, err := repositories.CheckGroupExistsByKey(group.Key, org_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Group already exists")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Group already exists"})
	}
	repositories.CreateGroup(&group)
	return c.JSON(http.StatusCreated, &group)

}

func DeleteGroup(c echo.Context) error {
	var group models.Group
	group_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	group_exists, group_err := repositories.CheckGroupExistsById(group_id)
	if group_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !group_exists {
		config.Log.Info("Group not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Group not exists"})
	}
	repositories.DeleteGroup(&group, group_id)
	return c.JSON(http.StatusNoContent, "")

}

func UpdateGroup(c echo.Context) error {
	var group models.Group
	var reqGroup models.Group
	group_id := string(c.Param("id"))
	org_id := string(c.Param("org_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Organization not exists"})
	}
	group_exists, group_err := repositories.CheckGroupExistsById(group_id)
	if group_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !group_exists {
		config.Log.Info("Group not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Group not exists"})
	}
	repositories.UpdateGroup(&group, &reqGroup, group_id)
	return c.JSON(http.StatusOK, &group)
}

func AddUserToGroup(c echo.Context) error {
	org_id := string(c.Param("org_id"))
	group_id := string(c.Param("id"))
	user_id := string(c.Param("user_id"))
	org_exists, org_err := repositories.CheckOrganizationExistsById(org_id)
	if org_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !org_exists {
		config.Log.Info("Organization not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Organization not exists"})
	}
	group_exists, group_err := repositories.CheckGroupExistsById(group_id)
	if group_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !group_exists {
		config.Log.Info("Group not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Group not exists"})
	}
	user_exists, user_err := repositories.CheckUserExistsById(user_id)
	if user_err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !user_exists {
		config.Log.Info("User not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User not exists"})
	}
	repositories.AddUserToGroup(group_id, user_id)
	return c.JSON(http.StatusOK, "")
}
