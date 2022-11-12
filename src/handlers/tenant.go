package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetTenants(org_id string, tenants *[]models.Tenant) error {
	return repositories.GetTenants(org_id, tenants)
}

func GetTenant(org_id string, id string, tenant *models.Tenant) error {
	return repositories.GetTenant(org_id, id, tenant)
}

func CreateTenant(org_id string, tenant *models.Tenant) error {
	return repositories.CreateTenant(org_id, tenant)
}

func DeleteTenant(org_id string, id string) error {
	return repositories.DeleteTenant(org_id, id)
}

func UpdateTenant(org_id string, id string, tenant *models.Tenant, reqTenant *models.Tenant) error {
	err := repositories.GetTenant(org_id, id, tenant)
	if err != nil {
		return err
	}
	tenant.Name = reqTenant.Name
	return repositories.UpdateTenant(tenant)
}

func CheckTenantExistsById(org_id string, id string) (bool, error) {
	var exists bool
	err := repositories.CheckTenantExistsById(org_id, id, &exists)
	if err != nil {
		return false, errors.New("tenant not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckTenantExistsByKey(org_id string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckTenantExistsByKey(org_id, key, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
