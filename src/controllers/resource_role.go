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

// @Description Get all resource roles.
// @Tags        Resource Role
// @Param res_id path int true "Resource ID"
// @Produce     json
// @Success     200 {array}  models.ResourceRole
// @failure     500
// @Router      /{res_id}/resource_role [get]
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
	err = handlers.GetResourceRoles(&resourceRoles, resId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRoles)
}

// @Description Get resource roles by ID.
// @Tags        Resource Role
// @Param res_id path int true "Resource ID"
// @Param id path int true "Resource Role ID"
// @Produce     json
// @Success     200 {object}  models.ResourceRoleWithGroupsUsers
// @failure     404,500
// @Router      /{res_id}/resource_role/{id} [get]
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
	err := handlers.GetUResourceRoleWithGroupsAndUsers(resroleId, &resourceRole)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)

}

// @Description Create resource role.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param request body models.ResourceRoleCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceRole
// @failure     400,403,500
// @Router      /{res_id}/resource_role [post]
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
	err = handlers.CreateResourceRoleAction(&resourceRole)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)

}

// @Description Delete resource role.
// @Tags        Resource Role
// @Param res_id path int true "Resource ID"
// @Param id path int true "Resource Role ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{res_id}/resource_role/{id} [delete]
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
		return utils.NotFoundErrorResponse("Resource Role")
	}
	err := handlers.DeleteResourceRole(&resourceRole, resRoleId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update resource role.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param id path int true "Resource Role ID"
// @Param request body models.ResourceRoleUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceRole
// @failure     400,403,500
// @Router      /{proj_id}/resource_role/{id} [put]
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
	err := handlers.UpdateResourceRole(&resourceRole, &reqResourceRole, resroleId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)
}

// @Description Assign resource role to user.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param user_id path int true "User ID"
// @Produce     json
// @Success     200
// @failure     403,404,500
// @Router      /{res_id}/resource_role/user/{user_id} [post]
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
	err = handlers.AddUserToResourceRole(resroleId, userId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, "")

}

// @Description Assign resource role to group.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param group_id path int true "Group ID"
// @Produce     json
// @Success     200
// @failure     403,404,500
// @Router      /{res_id}/resource_role/group/{group_id} [post]
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
	err = handlers.AddGroupToResourceRole(resroleId, groupId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, "")

}

// @Description Assign resource role to group.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path int true "Resource ID"
// @Param resact_id path int true "Resource Action ID"
// @Produce     json
// @Success     200
// @failure     403,404,500
// @Router      /{res_id}/resource_role/action/{resact_id} [post]
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
	err = handlers.AddResourceActionToResourceRole(resId, resroleId, resactId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, "")
}
