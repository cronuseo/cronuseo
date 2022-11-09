package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceActions(resourceActions *[]models.ResourceAction, res_id string) error {
	return config.DB.Model(&models.ResourceAction{}).Where("resource_id = ?", res_id).Find(&resourceActions).Error
}

func GetResourceAction(resourceAction *models.ResourceAction, resact_id string) error {
	return config.DB.Where("id = ?", resact_id).First(&resourceAction).Error
}

func CreateResourceAction(resourceAction *models.ResourceAction) error {
	return config.DB.Create(&resourceAction).Error
}

func DeleteResourceAction(resourceAction *models.ResourceAction, resact_id string) error {
	return config.DB.Where("id = ?", resact_id).Delete(&resourceAction).Error
}

func UpdateResourceAction(resourceAction *models.ResourceAction) error {
	return config.DB.Save(&resourceAction).Error
}

func DeleteAllResourceActions(res_id string) error {
	return config.DB.Where("resource_id = ?", res_id).Delete(&models.ResourceAction{}).Error
}

func CheckResourceActionExistsById(resactId string, exists *bool) error {
	return config.DB.Model(&models.ResourceAction{}).Select("count(*) > 0").Where("id = ?",
		resactId).Find(exists).Error
}

func CheckResourceActionExistsByKey(key string, resId string, exists *bool) error {
	return config.DB.Model(&models.ResourceAction{}).Select("count(*) > 0").Where(
		"key = ? AND resource_id = ?", key, resId).Find(exists).Error
}
