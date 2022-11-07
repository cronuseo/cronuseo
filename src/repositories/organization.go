package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetOrganizations(orgs *[]models.Organization) error {
	return config.DB.Find(&orgs).Error
}

func GetOrganization(org *models.Organization, id string) error {
	return config.DB.Where("id = ?", id).First(&org).Error
}

func CreateOrganization(org *models.Organization) error {
	return config.DB.Create(&org).Error
}

func DeleteOrganization(org *models.Organization, id string) error {
	return config.DB.Where("id = ?", id).Delete(&org).Error
}

func UpdateOrganization(org *models.Organization) error {
	return config.DB.Save(&org).Error
}

func CheckOrganizationExistsById(id string, exists *bool) error {
	return config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("id = ?", id).Find(exists).Error
}

func CheckOrganizationExistsByKey(key string, exists *bool) error {
	return config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("key = ?", key).Find(exists).Error
}
