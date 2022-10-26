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

func GetResourceRoles(c echo.Context) error {
	resourceRoles := []models.ResourceRole{}
	res_id := string(c.Param("res_id"))
	exists, err := repositories.CheckResourceExistsById(res_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	repositories.GetResourceRoles(&resourceRoles, res_id)
	return c.JSON(http.StatusOK, &resourceRoles)
}

func GetResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRoleWithGroupsUsers
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resrole_exists {
		config.Log.Info("Resource Action not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Action not exists"})
	}
	repositories.GetUResourceRoleWithGroupsAndUsers(resrole_id, &resourceRole)
	return c.JSON(http.StatusOK, &resourceRole)

}

func CreateResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	res_id := string(c.Param("res_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	if err := c.Bind(&resourceRole); err != nil {
		if resourceRole.Key == "" || len(resourceRole.Key) < 4 || resourceRole.Name == "" || len(resourceRole.Name) < 4 {
			return echo.NewHTTPError(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: err.Error()})
		}
	}
	if err := c.Validate(&resourceRole); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 400, Message: "Invalid inputs. Please check your inputs"})
	}
	int_res_id, _ := strconv.Atoi(res_id)
	resourceRole.ResourceID = int_res_id
	exists, err := repositories.CheckResourceRoleExistsByKey(resourceRole.Key, res_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if exists {
		config.Log.Info("Resource Role already exists")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Resource Role already exists"})
	}
	repositories.CreateResourceRoleAction(&resourceRole)
	return c.JSON(http.StatusOK, &resourceRole)

}

func DeleteResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	repositories.DeleteResourceRole(&resourceRole, resrole_id)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	var reqResourceRole models.ResourceRole
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	repositories.UpdateResourceRole(&resourceRole, &reqResourceRole, resrole_id)
	return c.JSON(http.StatusOK, &resourceRole)
}

func AddUserToResourceRole(c echo.Context) error {
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	user_id := string(c.Param("user_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	user_exists, user_err := repositories.CheckUserExistsById(user_id)
	if user_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !user_exists {
		config.Log.Info("User not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "User not exists"})
	}
	exists, err := repositories.CheckUserAlreadyAdded(resrole_id, user_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !exists {
		config.Log.Info("User already added")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "User already exists"})
	}
	repositories.AddUserToResourceRole(resrole_id, user_id)
	return c.JSON(http.StatusOK, "")

}

func AddGroupToResourceRole(c echo.Context) error {
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	group_id := string(c.Param("group_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	group_exists, group_err := repositories.CheckUserExistsById(group_id)
	if group_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !group_exists {
		config.Log.Info("User not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Group not exists"})
	}
	exists, err := repositories.CheckGroupAlreadyAdded(resrole_id, group_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !exists {
		config.Log.Info("Group already added")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Group already exists"})
	}
	repositories.AddGroupToResourceRole(resrole_id, group_id)
	return c.JSON(http.StatusOK, "")

}

func AddResourceActionToResourceRole(c echo.Context) error {
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	resact_id := string(c.Param("resact_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource not exists"})
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	resact_exists, resact_err := repositories.CheckResourceActionExistsById(resact_id)
	if resact_err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if !resact_exists {
		config.Log.Info("Resource Action not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	exists, err := repositories.CheckResourceActionAlreadyAdded(res_id, resrole_id, resact_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		return echo.NewHTTPError(http.StatusInternalServerError, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
	}
	if exists {
		config.Log.Info("Resource Action already added")
		return echo.NewHTTPError(http.StatusForbidden, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 403, Message: "Resource Action already exists"})
	}
	repositories.AddResourceActionToResourceRole(res_id, resrole_id, resact_id)
	return c.JSON(http.StatusOK, "")
}
