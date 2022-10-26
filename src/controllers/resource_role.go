package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/exceptions"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResourceRoles(c *gin.Context) {
	resourceRoles := []models.ResourceRole{}
	res_id := string(c.Param("res_id"))
	exists, err := repositories.CheckResourceExistsById(res_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	repositories.GetResourceRoles(&resourceRoles, res_id)
	c.JSON(http.StatusOK, &resourceRoles)
}

func GetResourceRoles2(c echo.Context) error {
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

func GetResourceRole(c *gin.Context) {
	var resourceRole models.ResourceRoleWithGroupsUsers
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resrole_exists {
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	repositories.GetUResourceRoleWithGroupsAndUsers(resrole_id, &resourceRole)
	c.JSON(http.StatusOK, &resourceRole)

}

func GetResourceRole2(c echo.Context) error {
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

func CreateResourceRole(c *gin.Context) {
	var resourceRole models.ResourceRole
	res_id := string(c.Param("res_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	if err := c.ShouldBindJSON(&resourceRole); err != nil {
		if resourceRole.Key == "" || len(resourceRole.Key) < 4 || resourceRole.Name == "" || len(resourceRole.Name) < 4 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: err.Error()})
			return
		}
	}
	int_res_id, _ := strconv.Atoi(res_id)
	resourceRole.ResourceID = int_res_id
	exists, err := repositories.CheckResourceRoleExistsByKey(resourceRole.Key, res_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Resource Role already exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role already exists"})
		return
	}
	repositories.CreateResourceRoleAction(&resourceRole)
	c.JSON(http.StatusOK, &resourceRole)

}

func CreateResourceRole2(c echo.Context) error {
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

func DeleteResourceRole(c *gin.Context) {
	var resourceRole models.ResourceRole
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	repositories.DeleteResourceRole(&resourceRole, resrole_id)
	c.JSON(http.StatusOK, "")
}

func DeleteResourceRole2(c echo.Context) error {
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

func UpdateResourceRole(c *gin.Context) {
	var resourceRole models.ResourceRole
	var reqResourceRole models.ResourceRole
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	repositories.UpdateResourceRole(&resourceRole, &reqResourceRole, resrole_id)
	c.JSON(http.StatusOK, &resourceRole)
}

func UpdateResourceRole2(c echo.Context) error {
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

func AddUserToResourceRole(c *gin.Context) {
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	user_id := string(c.Param("user_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	user_exists, user_err := repositories.CheckUserExistsById(user_id)
	if user_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !user_exists {
		config.Log.Info("User not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User not exists"})
		return
	}
	exists, err := repositories.CheckUserAlreadyAdded(resrole_id, user_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("User already added")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User already added"})
		return
	}
	repositories.AddUserToResourceRole(resrole_id, user_id)
	c.JSON(http.StatusOK, "")

}

func AddUserToResourceRole2(c echo.Context) error {
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

func AddGroupToResourceRole(c *gin.Context) {
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	group_id := string(c.Param("group_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	group_exists, group_err := repositories.CheckUserExistsById(group_id)
	if group_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !group_exists {
		config.Log.Info("User not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "User not exists"})
		return
	}
	exists, err := repositories.CheckGroupAlreadyAdded(resrole_id, group_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !exists {
		config.Log.Info("Group already added")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Group already added"})
		return
	}
	repositories.AddGroupToResourceRole(resrole_id, group_id)
	c.JSON(http.StatusOK, "")

}

func AddGroupToResourceRole2(c echo.Context) error {
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

func AddResourceActionToResourceRole(c *gin.Context) {
	res_id := string(c.Param("res_id"))
	resrole_id := string(c.Param("id"))
	resact_id := string(c.Param("resact_id"))
	res_exists, res_err := repositories.CheckResourceExistsById(res_id)
	if res_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !res_exists {
		config.Log.Info("Resource not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource not exists"})
		return
	}
	resrole_exists, resrole_err := repositories.CheckResourceRoleExistsById(resrole_id)
	if resrole_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resrole_exists {
		config.Log.Info("Resource Role not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Role not exists"})
		return
	}
	resact_exists, resact_err := repositories.CheckResourceActionExistsById(resact_id)
	if resact_err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if !resact_exists {
		config.Log.Info("Resource Action not exists")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action not exists"})
		return
	}
	exists, err := repositories.CheckResourceActionAlreadyAdded(res_id, resrole_id, resact_id)
	if err != nil {
		config.Log.Panic("Server Error!")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Server Error!"})
		return
	}
	if exists {
		config.Log.Info("Resource Action already added")
		c.AbortWithStatusJSON(http.StatusBadRequest, exceptions.Exception{Timestamp: time.Now().Format(time.RFC3339Nano), Status: 500, Message: "Resource Action already added"})
		return
	}
	repositories.AddResourceActionToResourceRole(res_id, resrole_id, resact_id)
	c.JSON(http.StatusOK, "")

}

func AddResourceActionToResourceRole2(c echo.Context) error {
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
