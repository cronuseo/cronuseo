package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetOrganizations(orgs *[]models.Organization) error {
	return repositories.GetOrganizations(orgs)
}

func GetOrganization(org *models.Organization, id string) error {
	return repositories.GetOrganization(org, id)
}

func CreateOrganization(org *models.Organization) error {
	return repositories.CreateOrganization(org)
}

func DeleteOrganization(id string) error {
	return repositories.DeleteOrganization(id)
}

func UpdateOrganization(org *models.Organization, reqOrg *models.Organization, id string) error {
	err := repositories.GetOrganization(org, id)
	if err != nil {
		return err
	}
	org.Name = reqOrg.Name
	return repositories.UpdateOrganization(org)
}

func CheckOrganizationExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckOrganizationExistsById(id, &exists)
	return exists, err
}

func CheckOrganizationExistsByKey(key string) (bool, error) {
	var exists bool
	err := repositories.CheckOrganizationExistsByKey(key, &exists)
	return exists, err

}
