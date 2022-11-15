package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all resource roles.
// @Tags        Resource Role
// @Param res_id path string true "Resource ID"
// @Produce     json
// @Success     200 {array}  models.ResourceRole
// @failure     500
// @Router      /{res_id}/resource_role [get]
func GetResourceRoles(c echo.Context) error {
	resourceRoles := []models.ResourceRole{}
	resId := string(c.Param("res_id"))

	exists, _ := handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	err := handlers.GetResourceRoles(resId, &resourceRoles)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRoles)
}

// @Description Get resource roles by ID.
// @Tags        Resource Role
// @Param res_id path string true "Resource ID"
// @Param id path string true "Resource Role ID"
// @Produce     json
// @Success     200 {object}  models.ResourceRole
// @failure     404,500
// @Router      /{res_id}/resource_role/{id} [get]
func GetResourceRole(c echo.Context) error {
	resId := string(c.Param("res_id"))
	var resourceRole models.ResourceRole
	resroleId := string(c.Param("id"))

	exists, _ := handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	exists, _ = handlers.CheckResourceRoleExistsById(resroleId)
	if !exists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}

	err := handlers.GetResourceRole(resId, resroleId, &resourceRole)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)

}

// @Description Create resource role.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path string true "Resource ID"
// @Param request body models.ResourceRoleCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceRole
// @failure     400,403,500
// @Router      /{res_id}/resource_role [post]
func CreateResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	resId := string(c.Param("res_id"))

	exists, _ := handlers.CheckResourceExistsById(resId)
	if !exists {
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

	resourceRole.ResourceID = resId
	exists, _ = handlers.CheckResourceRoleExistsByKey(resourceRole.Key, resId)
	if exists {
		config.Log.Info("Resource Role already exists")
		return utils.AlreadyExistsErrorResponse("Resource Role")
	}

	err := handlers.CreateResourceRole(resId, &resourceRole)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)

}

// @Description Delete resource role.
// @Tags        Resource Role
// @Param res_id path string true "Resource ID"
// @Param id path string true "Resource Role ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{res_id}/resource_role/{id} [delete]
func DeleteResourceRole(c echo.Context) error {

	resId := string(c.Param("res_id"))
	resRoleId := string(c.Param("id"))

	exists, _ := handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	exists, _ = handlers.CheckResourceRoleExistsById(resRoleId)
	if !exists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}
	err := handlers.DeleteResourceRole(resId, resRoleId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update resource role.
// @Tags        Resource Role
// @Accept      json
// @Param res_id path string true "Resource ID"
// @Param id path string true "Resource Role ID"
// @Param request body models.ResourceRoleUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceRole
// @failure     400,403,404,500
// @Router      /{res_id}/resource_role/{id} [put]
func UpdateResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	var reqResourceRole models.ResourceRoleUpdateRequest
	resId := string(c.Param("res_id"))
	resroleId := string(c.Param("id"))

	exists, _ := handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}
	exists, _ = handlers.CheckResourceRoleExistsById(resroleId)
	if !exists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}

	if err := c.Bind(&reqResourceRole); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&reqResourceRole); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.UpdateResourceRole(resId, resroleId, &resourceRole, &reqResourceRole)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)
}

// @Description Patch resource role.
// @Tags         Resource Role
// @Accept      json
// @Param res_id path string true "Resource ID"
// @Param id path string true "Resource Role ID"
// @Param request body models.ResourceRolePatchRequest true "body"
// @Produce     json
// @Success     201 {object}  models.ResourceRole
// @failure     400,403,404,500
// @Router      /{res_id}/resource_role/{id} [patch]
func PatchResourceRole(c echo.Context) error {
	var resourceRole models.ResourceRole
	var resourceRolePatch models.ResourceRolePatchRequest

	resId := string(c.Param("res_id"))
	resRoleId := string(c.Param("id"))

	exists, _ := handlers.CheckResourceExistsById(resId)
	if !exists {
		config.Log.Info("Resource not exists")
		return utils.NotFoundErrorResponse("Resource")
	}

	exists, _ = handlers.CheckResourceRoleExistsById(resRoleId)
	if !exists {
		config.Log.Info("Resource Role not exists")
		return utils.NotFoundErrorResponse("Resource Role")
	}

	if err := c.Bind(&resourceRolePatch); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&resourceRolePatch); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.PatchResourceRole(resId, resRoleId, &resourceRole, &resourceRolePatch)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &resourceRole)
}
