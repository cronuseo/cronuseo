package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description get all organizations.
// @Accept      */*
// @Produce     json
// @Success     200 {array}  models.Organization
// @failure     404 {string} Organization not exists "error"
// @failure     500 {string} string "Server Error!"
// @Router      /organization [get]
func GetOrganizations(c echo.Context) error {
	orgs := []models.Organization{}
	handlers.GetOrganizations(&orgs)
	return c.JSON(http.StatusOK, &orgs)
}

func GetOrganization(c echo.Context) error {
	var org models.Organization
	orgId := string(c.Param("id"))
	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	handlers.GetOrganization(&org, orgId)
	return c.JSON(http.StatusOK, &org)

}

func CreateOrganization(c echo.Context) error {
	var org models.Organization
	if err := c.Bind(&org); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&org); err != nil {
		return utils.InvalidErrorResponse()
	}
	exists, err := handlers.CheckOrganizationExistsByKey(org.Key)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Organization already exists")
		return utils.AlreadyExistsErrorResponse("Organization")
	}
	handlers.CreateOrganization(&org)
	return c.JSON(http.StatusCreated, &org)
}

func DeleteOrganization(c echo.Context) error {
	var org models.Organization
	orgId := string(c.Param("id"))
	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	handlers.DeleteOrganization(&org, orgId)
	return c.JSON(http.StatusNoContent, "")
}

func UpdateOrganization(c echo.Context) error {
	var org models.Organization
	orgId := string(c.Param("id"))
	var reqOrg models.Organization
	if err := c.Bind(&reqOrg); err != nil {
		if reqOrg.Name == "" || len(reqOrg.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}
	if err := c.Validate(&org); err != nil {
		return utils.InvalidErrorResponse()
	}
	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}
	handlers.UpdateOrganization(&org, &reqOrg, orgId)
	return c.JSON(http.StatusCreated, &org)
}
