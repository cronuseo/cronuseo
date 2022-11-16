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
// @Param org_id path string true "Organization ID"
// @Produce     json
// @Success     200 {array}  models.User
// @failure     500
// @Router      /{org_id}/user [get]
func GetUsers(c echo.Context) error {
	users := []models.User{}
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	err := handlers.GetUsers(orgId, &users)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &users)
}

// @Description Get user by ID.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Produce     json
// @Success     200 {object}  models.UserWithGroup
// @failure     404,500
// @Router      /{org_id}/user/{id} [get]
func GetUser(c echo.Context) error {
	var user models.User
	orgId := string(c.Param("org_id"))
	userId := string(c.Param("id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckUserExistsById(userId)
	if !exists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}

	err := handlers.GetUser(orgId, userId, &user)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &user)
}

// @Description Create user.
// @Tags        User
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body models.UserCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.User
// @failure     400,403,500
// @Router      /{org_id}/user [post]
func CreateUser(c echo.Context) error {
	var user models.User
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
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

	exists, _ = handlers.CheckUserExistsByUsername(orgId, user.Username)
	if exists {
		config.Log.Info("User already exists")
		return utils.AlreadyExistsErrorResponse("User")
	}

	err := handlers.CreateUser(orgId, &user)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &user)
}

// @Description Delete user.
// @Tags        User
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/user/{id} [delete]
func DeleteUser(c echo.Context) error {

	userId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckUserExistsById(userId)
	if !exists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}

	err := handlers.DeleteUser(orgId, userId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update user.
// @Tags        User
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "User ID"
// @Param request body models.UserUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.User
// @failure     400,403,404,500
// @Router      /{org_id}/user/{id} [put]
func UpdateUser(c echo.Context) error {
	var user models.User
	var reqUser models.UserUpdateRequest
	userId := string(c.Param("id"))
	orgId := string(c.Param("org_id"))

	exists, _ := handlers.CheckOrganizationExistsById(orgId)
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, _ = handlers.CheckUserExistsById(userId)
	if !exists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}

	if err := c.Bind(&reqUser); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&reqUser); err != nil {
		return utils.InvalidErrorResponse()
	}

	err := handlers.UpdateUser(orgId, userId, &user, &reqUser)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &user)
}
