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

func GetGroups(c echo.Context) error {
	groups := []models.Group{}
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
	handlers.GetGroups(&groups, orgId)
	return c.JSON(http.StatusOK, &groups)
}

func GetGroup(c echo.Context) error {
	var group models.GroupUsers
	orgId := string(c.Param("org_id"))
	groupId := string(c.Param("id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	groupExists, groupErr := handlers.CheckGroupExistsById(groupId)
	if groupErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !groupExists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}
	handlers.GetUsersFromGroup(groupId, &group)
	return c.JSON(http.StatusOK, &group)

}

func CreateGroup(c echo.Context) error {
	var group models.Group
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
	if err := c.Bind(&group); err != nil {
		if group.Key == "" || len(group.Key) < 4 || group.Name == "" || len(group.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&group); err != nil {
		return utils.InvalidErrorResponse()
	}
	intOrgId, _ := strconv.Atoi(orgId)
	group.OrganizationID = intOrgId
	exists, err := handlers.CheckGroupExistsByKey(group.Key, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Group already exists")
		return utils.AlreadyExistsErrorResponse("Group")
	}
	handlers.CreateGroup(&group)
	return c.JSON(http.StatusCreated, &group)

}

func DeleteGroup(c echo.Context) error {
	var group models.Group
	groupId := string(c.Param("id"))
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
	groupExists, groupErr := handlers.CheckGroupExistsById(groupId)
	if groupErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !groupExists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}
	handlers.DeleteGroup(&group, groupId)
	return c.JSON(http.StatusNoContent, "")

}

func UpdateGroup(c echo.Context) error {
	var group models.Group
	var reqGroup models.Group
	groupId := string(c.Param("id"))
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
	groupExists, groupErr := handlers.CheckGroupExistsById(groupId)
	if groupErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !groupExists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}
	handlers.UpdateGroup(&group, &reqGroup, groupId)
	return c.JSON(http.StatusOK, &group)
}

func AddUserToGroup(c echo.Context) error {
	orgId := string(c.Param("org_id"))
	groupId := string(c.Param("id"))
	userId := string(c.Param("user_id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	groupExists, groupErr := handlers.CheckGroupExistsById(groupId)
	if groupErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !groupExists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
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
	handlers.AddUserToGroup(groupId, userId)
	return c.JSON(http.StatusOK, "")
}
