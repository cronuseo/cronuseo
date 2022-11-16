package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all groups.
// @Tags        Group
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  models.Group
// @failure     500
// @Router      /{org_id}/group [get]
func GetGroups(c echo.Context) error {
	groups := []models.Group{}
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	err := handlers.GetGroups(orgId, &groups)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &groups)
}

// @Description Get groups by ID.
// @Tags        Group
// @Param org_id path string true "Organization ID"
// @Param id path string true "Group ID"
// @Produce     json
// @Success     200 {object}  models.Group
// @failure     404,500
// @Router      /{org_id}/group/{id} [get]
func GetGroup(c echo.Context) error {
	var group models.Group
	orgId := string(c.Param("org_id"))
	groupId := string(c.Param("id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckGroupExistsById(groupId)
	if !exists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}

	err := handlers.GetGroup(orgId, groupId, &group)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)

}

// @Description Create group.
// @Tags        Group
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body models.GroupCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,500
// @Router      /{org_id}/group [post]
func CreateGroup(c echo.Context) error {
	var group models.Group
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
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

	group.OrgID = orgId
	exists, _ = handlers.CheckGroupExistsByKey(group.Key, orgId)
	if exists {
		config.Log.Info("Group already exists")
		return utils.AlreadyExistsErrorResponse("Group")
	}

	err := handlers.CreateGroup(orgId, &group)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &group)

}

// @Description Delete group.
// @Tags        Group
// @Param org_id path string true "Organization ID"
// @Param id path string true "Group ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/group/{id} [delete]
func DeleteGroup(c echo.Context) error {

	groupId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckGroupExistsById(groupId)
	if !exists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}

	err := handlers.DeleteGroup(orgId, groupId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")

}

// @Description Update group.
// @Tags        Group
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Group ID"
// @Param request body models.GroupUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,404,500
// @Router      /{org_id}/group/{id} [put]
func UpdateGroup(c echo.Context) error {
	var group models.Group
	var reqGroup models.GroupUpdateRequest
	groupId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckGroupExistsById(groupId)
	if !exists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}

	if err := c.Bind(&reqGroup); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&reqGroup); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.UpdateGroup(orgId, groupId, &group, &reqGroup)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)
}

// @Description Patch group.
// @Tags        Group
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Group ID"
// @Param request body models.GroupPatchRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,404,500
// @Router      /{org_id}/group/{id} [patch]
func PatchGroup(c echo.Context) error {
	var group models.Group
	var groupPatch models.GroupPatchRequest
	groupId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckGroupExistsById(groupId)
	if !exists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}

	if err := c.Bind(&groupPatch); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&groupPatch); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.PatchGroup(orgId, groupId, &group, &groupPatch)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)
}
