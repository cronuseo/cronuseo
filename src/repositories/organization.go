package repositories

import (
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

func UpdateOrganization(org *models.Organization) {
	config.DB.Save(&org)
}

func CheckOrganizationExistsById(id string, exists *bool) error {
	return config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("id = ?", id).Find(exists).Error
}

func CheckOrganizationExistsByKey(key string, exists *bool) error {
	return config.DB.Model(&models.Organization{}).Select("count(*) > 0").Where("key = ?", key).Find(exists).Error
}
