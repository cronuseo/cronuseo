package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/handlers"
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/utils"
)

// @Description Get all tenants.
// @Tags        Tenant
// @Produce     json
// @Param org_id path string true "Organization ID"
// @Success     200 {array}  models.Tenant
// @failure     500
// @Router      /{org_id}/tenant [get]
func GetTenants(c echo.Context) error {
	tenant := []models.Tenant{}
	orgId := string(c.Param("org_id"))
	err := handlers.GetTenants(orgId, &tenant)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &tenant)
}

// @Description Get tenant by ID.
// @Tags        Tenant
// @Param org_id path string true "Organization ID"
// @Param id path string true "Tenant ID"
// @Produce     json
// @Success     200 {object}  models.Tenant
// @failure     404,500
// @Router      /{org_id}/tenant/{id} [get]
func GetTenant(c echo.Context) error {
	var tenant models.Tenant
	orgId := string(c.Param("org_id"))
	tenantId := string(c.Param("id"))

	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, err = handlers.CheckTenantExistsById(orgId, tenantId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	err = handlers.GetTenant(orgId, tenantId, &tenant)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusOK, &tenant)

}

// @Description Create tenant.
// @Tags        Tenant
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param request body models.TenantCreateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Tenant
// @failure     400,403,500
// @Router      /{org_id}/tenant [post]
func CreateTenant(c echo.Context) error {
	var tenant models.Tenant
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

	if err := c.Bind(&tenant); err != nil {
		return utils.InvalidErrorResponse()
	}
	if err := c.Validate(&tenant); err != nil {
		return utils.InvalidErrorResponse()
	}

	exists, err = handlers.CheckTenantExistsByKey(orgId, tenant.Key)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if exists {
		config.Log.Info("Tenant already exists")
		return utils.AlreadyExistsErrorResponse("Organization")
	}

	err = handlers.CreateTenant(orgId, &tenant)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &tenant)
}

// @Description Delete tenant.
// @Tags        Tenant
// @Param org_id path string true "Organization ID"
// @Param id path string true "Tenant ID"
// @Produce     json
// @Success     204
// @failure     404,500
// @Router      /{org_id}/tenant/{id} [delete]
func DeleteTenant(c echo.Context) error {

	orgId := string(c.Param("org_id"))
	tenantId := string(c.Param("id"))

	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, err = handlers.CheckTenantExistsById(orgId, tenantId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	err = handlers.DeleteTenant(orgId, tenantId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusNoContent, "")
}

// @Description Update tenant.
// @Tags        Tenant
// @Accept      json
// @Param org_id path string true "Organization ID"
// @Param id path string true "Tenant ID"
// @Param request body models.TenantUpdateRequest true "body"
// @Produce     json
// @Success     201 {object}  models.Tenant
// @failure     400,403,404,500
// @Router      /{org_id}/tenant/{id} [put]
func UpdateTenant(c echo.Context) error {
	var tenant models.Tenant
	orgId := string(c.Param("org_id"))
	tenantId := string(c.Param("id"))
	var reqTenant models.Tenant

	exists, err := handlers.CheckOrganizationExistsById(orgId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Organization not exists")
		return utils.NotFoundErrorResponse("Organization")
	}

	exists, err = handlers.CheckTenantExistsById(orgId, tenantId)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	if !exists {
		config.Log.Info("Tenant not exists")
		return utils.NotFoundErrorResponse("Tenant")
	}

	if err := c.Bind(&reqTenant); err != nil {
		if reqTenant.Name == "" || len(reqTenant.Name) < 4 {
			return utils.InvalidErrorResponse()
		}
	}

	err = handlers.UpdateTenant(orgId, tenantId, &tenant, &reqTenant)
	if err != nil {
		config.Log.Panic("Server Error!")
		return utils.ServerErrorResponse()
	}
	return c.JSON(http.StatusCreated, &tenant)
}
