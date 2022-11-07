package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all organizations.
// @Tags        Organization
// @Produce     json
// @Success     200 {array}  models.Organization
// @failure     500
// @Router      /organization [get]
func GetOrganizations(c echo.Context) error {
	orgs := []models.Organization{}
	err := handlers.GetOrganizations(&orgs)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &orgs)
}

// @Description Get organization by ID.
// @Tags        Organization
// @Param org_id path int true "Organization ID"
// @Produce     json
// @Success     200 {object}  models.Organization
// @failure     404,500
// @Router      /organization/{org_id} [get]
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
	err = handlers.GetOrganization(&org, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &org)

}

// @Description Create organization.
// @Tags        Organization
// @Accept      json
// @Param request body models.OrganizationRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Organization
// @failure     400,403,500
// @Router      /organization [post]
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
	err = handlers.CreateOrganization(&org)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &org)
}

// @Description Delete organization.
// @Tags        Organization
// @Param org_id path int true "Organization ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /organization/{org_id} [delete]
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
	err = handlers.DeleteOrganization(&org, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update organization.
// @Tags        Organization
// @Accept      json
// @Param org_id path int true "Organization ID"
// @Param request body models.OrganizationUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Organization
// @failure     400,403,500
// @Router      /organization/{org_id} [put]
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
	err = handlers.UpdateOrganization(&org, &reqOrg, orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &org)
}
