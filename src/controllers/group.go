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
// @Param tenant_id path string true "Tenant ID"
// @Produce     json
// @Success     200 {array}  models.Group
// @failure     500
// @Router      /{tenant_id}/group [get]
func GetGroups(c echo.Context) error {
	groups := []models.Group{}
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	err := handlers.GetGroups(tenantId, &groups)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &groups)
}

// @Description Get groups by ID.
// @Tags        Group
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Group ID"
// @Produce     json
// @Success     200 {object}  models.Group
// @failure     404,500
// @Router      /{tenant_id}/group/{id} [get]
func GetGroup(c echo.Context) error {
	var group models.Group
	tenantId := string(c.Param("tenant_id"))
	groupId := string(c.Param("id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	exists, _ = handlers.CheckGroupExistsById(groupId)
	if !exists {
		config.Log.Info("Group not exists")
		return utils.NotFoundErrorResponse("Group")
	}

	err := handlers.GetGroup(tenantId, groupId, &group)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)

}

// @Description Create group.
// @Tags        Group
// @Accept      json
// @Param tenant_id path string true "Tenant ID"
// @Param request body models.GroupCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,500
// @Router      /{tenant_id}/group [post]
func CreateGroup(c echo.Context) error {
	var group models.Group
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	if err := c.Bind(&group); err != nil {
		if group.Key == "" || len(group.Key) < 4 || group.Name == "" || len(group.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&group); err != nil {
		return utils.InvalidErrorResponse()
	}

	group.TenantID = tenantId
	exists, _ = handlers.CheckGroupExistsByKey(group.Key, tenantId)
	if exists {
		config.Log.Info("Group already exists")
		return utils.AlreadyExistsErrorResponse("Group")
	}

	err := handlers.CreateGroup(tenantId, &group)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &group)

}

// @Description Delete group.
// @Tags        Group
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Group ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{tenant_id}/group/{id} [delete]
func DeleteGroup(c echo.Context) error {
	var group models.Group
	groupId := string(c.Param("id"))
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	exists, _ = handlers.CheckGroupExistsById(groupId)
	if !exists {
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

	err := handlers.DeleteGroup(tenantId, groupId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")

}

// @Description Update group.
// @Tags        Group
// @Accept      json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Group ID"
// @Param request body models.GroupUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Group
// @failure     400,403,404,500
// @Router      /{tenant_id}/group/{id} [put]
func UpdateGroup(c echo.Context) error {
	var group models.Group
	var reqGroup models.GroupUpdateRequest
	groupId := string(c.Param("id"))
	tenantId := string(c.Param("tenant_id"))

	exists, _ := handlers.CheckTenantExistsById(tenantId)
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
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

	err := handlers.UpdateGroup(tenantId, groupId, &group, &reqGroup)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &group)
}

// @Description Add users to group.
// @Tags        Group
// @Accept      json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Group ID"
// @Param request body models.AddUsersToGroup true "body"
// @Produce     json
// @Success     200
// @failure     404,500
// @Router      /{tenant_id}/group/{id}/user [post]
// func AddUsersToGroup(c echo.Context) error {
// 	tenantId := string(c.Param("tenant_id"))
// 	groupId := string(c.Param("id"))
// 	users := models.AddUsersToGroup{}
// 	orgExists, orgErr := handlers.CheckTenantExistsById(tenantId)
// 	if orgErr != nil {
// 		config.Log.Panic("Server Error!")
// 		return utils.ServerErrorResponse()
// 	}
// 	if !orgExists {
// 		config.Log.Info("Tenant not exists")
// 		return utils.NotFoundErrorResponse("Tenant")
// 	}
// 	groupExists, groupErr := handlers.CheckGroupExistsById(groupId)
// 	if groupErr != nil {
// 		config.Log.Panic("Server Error!")
// 		return utils.ServerErrorResponse()
// 	}
// 	if !groupExists {
// 		config.Log.Info("Group not exists")
// 		return utils.NotFoundErrorResponse("Group")
// 	}
// 	if err := c.Bind(&users); err != nil {
// 		return utils.InvalidErrorResponse()
// 	}
// 	if err := c.Validate(&users); err != nil {
// 		return utils.InvalidErrorResponse()
// 	}
// 	userExists, userErr := handlers.CheckAllUsersExistsById(users)
// 	if userErr != nil {
// 		config.Log.Panic("Server Error!")
// 		return utils.ServerErrorResponse()
// 	}
// 	if !userExists {
// 		config.Log.Info("Users not exists")
// 		return utils.NotFoundErrorResponse("Users")
// 	}

// 	err := handlers.AddUsersToGroup(groupId, users)
// 	if err != nil {
// 		config.Log.Panic("Server Error!")
// 		return utils.ServerErrorResponse()
// 	}
// 	return c.JSON(http.StatusOK, "")
// }
