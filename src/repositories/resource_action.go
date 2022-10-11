package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResourceActions(resourceActions *[]models.ResourceAction, res_id string) {
	config.DB.Model(&models.ResourceAction{}).Where("resource_id = ?", res_id).Find(&resourceActions)
}

func GetResourceAction(resourceAction *models.ResourceAction, resact_id string) {
	config.DB.Where("id = ?", resact_id).First(&resourceAction)
}

func CreateResourceActionAction(resourceAction *models.ResourceAction) {
	config.DB.Create(&resourceAction)
}

func DeleteResourceAction(resourceAction *models.ResourceAction, resact_id string) {
	config.DB.Where("id = ?", resact_id).Delete(&resourceAction)
}

func UpdateResourceAction(resourceAction *models.ResourceAction, reqResourceAction *models.ResourceAction, resact_id string) {
	config.DB.Where("id = ?", resact_id).First(&resourceAction)
	resourceAction.Name = reqResourceAction.Name
	resourceAction.Description = reqResourceAction.Description
	config.DB.Save(&resourceAction)
}

func CheckResourceActionExistsById(resact_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceAction{}).Select("count(*) > 0").Where("id = ?", resact_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("resource action not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckResourceActionExistsByKey(key string, res_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.ResourceAction{}).Select("count(*) > 0").Where("key = ? AND resource_id = ?", key, res_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
