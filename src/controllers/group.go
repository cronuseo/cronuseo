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

// @Description Get all groups.
// @Tags        Group
// @Param org_id path int true "Organization ID"
// @Produce     json
// @Success     200 {array}  models.Group
// @failure     500
// @Router      /{org_id}/group [get]
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
	err = handlers.GetGroups(&groups, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &groups)
}

// @Description Get groups by ID.
// @Tags        Group
// @Param org_id path int true "Organization ID"
// @Param id path int true "Group ID"
// @Produce     json
// @Success     200 {object}  models.GroupUsers
// @failure     404,500
// @Router      /{org_id}/group/{id} [get]
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
	err := handlers.GetUsersFromGroup(groupId, &group)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)

}

// @Description Create group.
// @Tags        Group
// @Accept      json
// @Param org_id path int true "Organization ID"
// @Param request body models.GroupCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,500
// @Router      /{org_id}/group [post]
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
	err = handlers.CreateGroup(&group)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &group)

}

// @Description Delete group.
// @Tags        Group
// @Param org_id path int true "Organization ID"
// @Param id path int true "Group ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/group/{id} [delete]
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
	if err := c.Bind(&group); err != nil {
		if group.Name == "" || len(group.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&group); err != nil {
		return utils.InvalidErrorResponse()
	}
	err := handlers.DeleteGroup(&group, groupId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")

}

// @Description Update group.
// @Tags        Group
// @Accept      json
// @Param org_id path int true "Organization ID"
// @Param id path int true "Group ID"
// @Param request body models.GroupUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,404,500
// @Router      /{org_id}/group/{id} [put]
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
	err := handlers.UpdateGroup(&group, &reqGroup, groupId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)
}

// @Description Add users to group.
// @Tags        Group
// @Accept      json
// @Param org_id path int true "Organization ID"
// @Param id path int true "Group ID"
// @Param request body models.AddUsersToGroup true "body"
// @Produce     json
// @Success     200
// @failure     404,500
// @Router      /{org_id}/group/{id}/user [post]
func AddUsersToGroup(c echo.Context) error {
	orgId := string(c.Param("org_id"))
	groupId := string(c.Param("id"))
	users := models.AddUsersToGroup{}
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
	if err := c.Bind(&users); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&users); err != nil {
		return utils.InvalidErrorResponse()
	}
	userExists, userErr := handlers.CheckAllUsersExistsById(users)
	if userErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !userExists {
		config.Log.Info("Users not exists")
		return utils.NotFoundErrorResponse("Users")
	}

	err := handlers.AddUsersToGroup(groupId, users)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, "")
}
