package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetOrganizations(orgs *[]models.Organization) {
	config.DB.Find(&orgs)
}

func GetOrganization(org *models.Organization, id string) {
	config.DB.Where("id = ?", id).First(&org)
}

func CreateOrganization(org *models.Organization) {
	config.DB.Create(&org)
}

func DeleteOrganization(org *models.Organization, id string) {
	config.DB.Where("id = ?", id).Delete(&org)
}

func UpdateOrganization(org *models.Organization, reqOrg *models.Organization, id string) {
	config.DB.Where("id = ?", id).First(&org)
	org.Name = reqOrg.Name
	config.DB.Save(&org)
}

func CheckOrganizationExistsById(id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error
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
	err := config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("key = ?", key).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
