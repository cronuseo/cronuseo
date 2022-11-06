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

func GetUsers(c echo.Context) error {
	users := []models.User{}
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
	handlers.GetUsers(&users, orgId)
	return c.JSON(http.StatusOK, &users)
}

func GetUser(c echo.Context) error {
	var user models.UserWithGroup
	orgId := string(c.Param("org_id"))
	userId := string(c.Param("id"))
	orgExists, orgErr := handlers.CheckOrganizationExistsById(orgId)
	if orgErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !orgExists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
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
	handlers.GetUser(&user, userId)
	return c.JSON(http.StatusOK, &user)
}

func CreateUser(c echo.Context) error {
	var user models.User
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
	if err := c.Bind(&user); err != nil {
		if user.Username == "" || len(user.Username) < 4 || user.FirstName == "" ||
			len(user.FirstName) < 4 || user.LastName == "" || len(user.LastName) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&user); err != nil {
		return utils.InvalidErrorResponse()
	}
	intOrgId, _ := strconv.Atoi(orgId)
	user.OrganizationID = intOrgId
	exists, err := handlers.CheckUserExistsByUsername(user.Username, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("User already exists")
		return utils.AlreadyExistsErrorResponse("User")
	}
	handlers.CreateUser(&user)
	return c.JSON(http.StatusCreated, &user)
}

func DeleteUser(c echo.Context) error {
	var user models.User
	userId := string(c.Param("id"))
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
	userExists, userErr := handlers.CheckUserExistsById(userId)
	if userErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !userExists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}
	handlers.DeleteUser(&user, userId)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateUser(c echo.Context) error {
	var user models.User
	var reqUser models.User
	userId := string(c.Param("id"))
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
	userExists, userErr := handlers.CheckUserExistsById(userId)
	if userErr != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !userExists {
		config.Log.Info("User not exists")
		return utils.NotFoundErrorResponse("User")
	}
	handlers.UpdateUser(&user, &reqUser, userId)
	return c.JSON(http.StatusCreated, &user)
}
