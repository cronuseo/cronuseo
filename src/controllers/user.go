package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all users.
// @Tags        User
// @Param tenant_id path int true "Tenant ID"
// @Produce     json
// @Success     200 {array}  models.User
// @failure     500
// @Router      /{tenant_id}/user [get]
func GetUsers(c echo.Context) error {
	users := []models.User{}
	tenantId := string(c.Param("tenant_id"))
	exists, err := handlers.CheckTenantExistsById(tenantId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}
	err = handlers.GetUsers(tenantId, &users)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &users)
}

// @Description Get user by ID.
// @Tags        User
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "User ID"
// @Produce     json
// @Success     200 {object}  models.UserWithGroup
// @failure     404,500
// @Router      /{tenant_id}/user/{id} [get]
func GetUser(c echo.Context) error {
	var user models.User
	tenantId := string(c.Param("tenant_id"))
	userId := string(c.Param("id"))

	orgExists, orgErr := handlers.CheckTenantExistsById(tenantId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
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

	err := handlers.GetUser(tenantId, userId, &user)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &user)
}

// @Description Create user.
// @Tags        User
// @Accept      json
// @Param tenant_id path int true "Tenant ID"
// @Param request body models.UserCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.User
// @failure     400,403,500
// @Router      /{tenant_id}/user [post]
func CreateUser(c echo.Context) error {
	var user models.User
	tenantId := string(c.Param("tenant_id"))

	orgExists, orgErr := handlers.CheckTenantExistsById(tenantId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	if err := c.Bind(&user); err != nil {
		if user.Username == "" || len(user.Username) < 4 || user.FirstName == "" ||
			len(user.FirstName) < 4 || user.LastName == "" || len(user.LastName) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&user); err != nil {
		return utils.InvalidErrorResponse()
	}

	exists, err := handlers.CheckUserExistsByUsername(tenantId, user.Username)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("User already exists")
		return utils.AlreadyExistsErrorResponse("User")
	}

	err = handlers.CreateUser(tenantId, &user)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &user)
}

// @Description Delete user.
// @Tags        User
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "User ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{tenant_id}/user/{id} [delete]
func DeleteUser(c echo.Context) error {

	userId := string(c.Param("id"))
	tenantId := string(c.Param("tenant_id"))

	orgExists, orgErr := handlers.CheckTenantExistsById(tenantId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
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

	err := handlers.DeleteUser(tenantId, userId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update user.
// @Tags        User
// @Accept      json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "User ID"
// @Param request body models.UserUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.User
// @failure     400,403,404,500
// @Router      /{tenant_id}/user/{id} [put]
func UpdateUser(c echo.Context) error {
	var user models.User
	var reqUser models.User
	userId := string(c.Param("id"))
	tenantId := string(c.Param("tenant_id"))

	orgExists, orgErr := handlers.CheckTenantExistsById(tenantId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
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

	err := handlers.UpdateUser(tenantId, userId, &user, &reqUser)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &user)
}
