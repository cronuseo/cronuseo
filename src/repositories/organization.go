package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetOrganizations(orgnizations *[]models.Organization) error {
	err := config.DB.Select(orgnizations, "SELECT * FROM organization")
	if err != nil {
		return err
	}
	return nil
}

func GetOrganization(org *models.Organization, id string) error {
	err := config.DB.Get(org, "SELECT * FROM organization WHERE org_id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func CreateOrganization(org *models.Organization) error {
	stmt, err := config.DB.Prepare("INSERT INTO organization(org_key,name) VALUES($1, $2)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(org.Key, org.Name)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOrganization(id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM organization WHERE org_id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateOrganization(org *models.Organization) error {
	stmt, err := config.DB.Prepare("UPDATE organization SET name = $1 WHERE org_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(org.Name, org.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckOrganizationExistsById(id string, exists *bool) error {
	return config.DB.QueryRow("SELECT exists (SELECT org_id FROM organization WHERE org_id = $1)", id).Scan(exists)
}

func CheckOrganizationExistsByKey(key string, exists *bool) error {
	return config.DB.QueryRow("SELECT exists (SELECT org_id FROM organization WHERE org_key = $1)", key).Scan(exists)
}
