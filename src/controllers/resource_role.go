package controllers

import (
	"github.com/shashimalcse/Cronuseo/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(c echo.Context) error {
	var resourceRoles []models.ResourceRole
	resId := string(c.Param("res_id"))
	exists, err := handlers.CheckResourceExistsById(resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	handlers.GetResourceRoles(&resourceRoles, resId)
	return c.JSON(http.StatusOK, &resourceRoles)
}

func GetResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRoleWithGroupsUsers
	resId := string(c.Param("res_id"))
	resroleId := string(c.Param("id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resroleExists, resroleErr := handlers.CheckResourceRoleExistsById(resroleId)
	if resroleErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resroleExists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}
	handlers.GetUResourceRoleWithGroupsAndUsers(resroleId, &resourceRole)
	return c.JSON(http.StatusOK, &resourceRole)

}

func CreateResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	resId := string(c.Param("res_id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	if err := c.Bind(&resourceRole); err != nil {
		if resourceRole.Key == "" || len(resourceRole.Key) < 4 || resourceRole.Name == "" || len(resourceRole.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&resourceRole); err != nil {
		return utils.InvalidErrorResponse()
	}
	intResId, _ := strconv.Atoi(resId)
	resourceRole.ResourceID = intResId
	exists, err := handlers.CheckResourceRoleExistsByKey(resourceRole.Key, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Resource Role already exists")
		return utils.AlreadyExistsErrorResponse("Resource Role")
	}
	handlers.CreateResourceRoleAction(&resourceRole)
	return c.JSON(http.StatusOK, &resourceRole)

}

func DeleteResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	resId := string(c.Param("res_id"))
	resRoleId := string(c.Param("id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resroleExists, resroleErr := handlers.CheckResourceRoleExistsById(resRoleId)
	if resroleErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resroleExists {
		config.Log.Info("Resource Role not exists")
		return echo.NewHTTPError(http.StatusNotFound, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 404, Message: "Resource Role not exists"})
	}
	handlers.DeleteResourceRole(&resourceRole, resRoleId)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	var reqResourceRole models.ResourceRole
	resId := string(c.Param("res_id"))
	resroleId := string(c.Param("id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resroleExists, resroleErr := handlers.CheckResourceRoleExistsById(resroleId)
	if resroleErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resroleExists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}
	handlers.UpdateResourceRole(&resourceRole, &reqResourceRole, resroleId)
	return c.JSON(http.StatusOK, &resourceRole)
}

func AddUserToResourceRole(c echo.Context) error {
	resId := string(c.Param("res_id"))
	resroleId := string(c.Param("id"))
	userId := string(c.Param("user_id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resroleExists, resroleErr := handlers.CheckResourceRoleExistsById(resroleId)
	if resroleErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resroleExists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}
	userExists, userErr := handlers.CheckUserExistsById(userId)
	if userErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !userExists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}
	exists, err := handlers.CheckUserAlreadyAdded(resroleId, userId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("User already added")
		return utils.AlreadyExistsErrorResponse("User")
	}
	handlers.AddUserToResourceRole(resroleId, userId)
	return c.JSON(http.StatusOK, "")

}

func AddGroupToResourceRole(c echo.Context) error {
	resId := string(c.Param("res_id"))
	resroleId := string(c.Param("id"))
	groupId := string(c.Param("group_id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resroleExists, resroleErr := handlers.CheckResourceRoleExistsById(resroleId)
	if resroleErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resroleExists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}
	groupExists, groupErr := handlers.CheckUserExistsById(groupId)
	if groupErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !groupExists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}
	exists, err := handlers.CheckGroupAlreadyAdded(resroleId, groupId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Group already added")
		return utils.AlreadyExistsErrorResponse("Group")
	}
	handlers.AddGroupToResourceRole(resroleId, groupId)
	return c.JSON(http.StatusOK, "")

}

func AddResourceActionToResourceRole(c echo.Context) error {
	resId := string(c.Param("res_id"))
	resroleId := string(c.Param("id"))
	resactId := string(c.Param("resact_id"))
	resExists, resErr := handlers.CheckResourceExistsById(resId)
	if resErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resExists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	resroleExists, resroleErr := handlers.CheckResourceRoleExistsById(resroleId)
	if resroleErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resroleExists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}
	resactExists, resactErr := handlers.CheckResourceActionExistsById(resactId)
	if resactErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !resactExists {
		config.Log.Info("Resource Action not exists")
		return utils.NotFoundErrorResponse("Resource Action")
	}
	exists, err := handlers.CheckResourceActionAlreadyAdded(resId, resroleId, resactId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Resource Action already added")
		return utils.AlreadyExistsErrorResponse("Resource Action")
	}
	handlers.AddResourceActionToResourceRole(resId, resroleId, resactId)
	return c.JSON(http.StatusOK, "")
}
