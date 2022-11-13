package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetTenants(org_id string, tenants *[]models.Tenant) error {
	err := config.DB.Select(tenants, "SELECT * FROM tenant WHERE org_id = $1", org_id)
	if err != nil {
		return err
	}
	return nil
}

func GetTenant(org_id string, id string, tenant *models.Tenant) error {
	err := config.DB.Get(tenant, "SELECT * FROM tenant WHERE org_id = $1 AND tenant_id = $2", org_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateTenant(org_id string, tenant *models.Tenant) error {
	stmt, err := config.DB.Prepare("INSERT INTO tenant(tenant_key,name,org_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tenant.Key, tenant.Name, org_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTenant(org_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM tenant WHERE org_id = $1 AND tenant_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(org_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTenant(tenant *models.Tenant) error {
	stmt, err := config.DB.Prepare("UPDATE tenant SET name = $1 WHERE org_id = $2 AND tenant_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tenant.Name, tenant.OraganizationID, tenant.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckTenantExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT tenant_id FROM tenant WHERE tenant_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckTenantExistsByKey(org_id string, key string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT tenant_key FROM tenant WHERE org_id = $1 AND tenant_key = $2)",
		org_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
