package handlers

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetOrganizations(orgs *[]models.Organization) {
	repositories.GetOrganizations(orgs)
}

func GetOrganization(org *models.Organization, id string) {
	repositories.GetOrganization(org, id)
}

func CreateOrganization(org *models.Organization) {
	repositories.CreateOrganization(org)
}

func DeleteOrganization(org *models.Organization, id string) {
	repositories.DeleteOrganization(org, id)
}

func UpdateOrganization(org *models.Organization, reqOrg *models.Organization, id string) {
	repositories.GetOrganization(org, id)
	org.Name = reqOrg.Name
	repositories.UpdateOrganization(org)
}

func CheckOrganizationExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckOrganizationExistsById(id, &exists)
	if err != nil {
		return false, errors.New("organization not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckOrganizationExistsByKey(key string) (bool, error) {
	var exists bool
	err := repositories.CheckOrganizationExistsByKey(key, &exists)
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
