package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceRoles(resourceRoles *[]models.ResourceRole, res_id string) {
	config.DB.Model(&models.ResourceRole{}).Where("resource_id = ?", res_id).Find(&resourceRoles)
}

func GetResourceRole(resourceRole *models.ResourceRole, resact_id string) {
	config.DB.Where("id = ?", resact_id).First(&resourceRole)
}

func CreateResourceRoleAction(resourceRole *models.ResourceRole) {
	config.DB.Create(&resourceRole)
}

func DeleteResourceRole(resourceRole *models.ResourceRole, resact_id string) {
	config.DB.Where("id = ?", resact_id).Delete(&resourceRole)
}

func UpdateResourceRole(resourceRole *models.ResourceRole, reqResourceRole *models.ResourceRole, resact_id string) {
	config.DB.Where("id = ?", resact_id).First(&resourceRole)
	resourceRole.Name = reqResourceRole.Name
	resourceRole.Description = reqResourceRole.Description
	config.DB.Save(&resourceRole)
}

func DeleteAllResourceRoles(res_id string) {
	config.DB.Where("resource_id = ?", res_id).Delete(&models.ResourceRole{})
}

func CheckResourceRoleExistsById(resact_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("id = ?", resact_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("resource role not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckResourceRoleExistsByKey(key string, res_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceRole{}).Select("count(*) > 0").Where("key = ? AND resource_id = ?", key, res_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
