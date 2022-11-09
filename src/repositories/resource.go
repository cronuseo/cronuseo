package repositories

import (
	"fmt"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResources(resources *[]models.Resource, proj_id string) error {
	return config.DB.Model(&models.Resource{}).Where("project_id = ?", proj_id).Find(&resources).Error
}

func GetResource(resource *models.Resource, res_id string) error {
	return config.DB.Where("id = ?", res_id).First(&resource).Error
}

func CreateResource(resource *models.Resource) error {
	return config.DB.Create(&resource).Error
}

func DeleteResource(resource *models.Resource, res_id string) error {
	return config.DB.Where("id = ?", res_id).Delete(&resource).Error
}

func UpdateResource(resource *models.Resource) error {
	return config.DB.Save(&resource).Error
}

func DeleteAllResources(proj_id string) error {
	resources := []models.Resource{}
	err := GetResources(&resources, proj_id)
	if err != nil {
		return err
	}
	for _, resource := range resources {
		res_id := resource.ID
		err = DeleteAllResourceActions(fmt.Sprint(res_id))
		if err != nil {
			return err
		}
		err = DeleteAllResourceRoles(fmt.Sprint(res_id))
		if err != nil {
			return err
		}
	}
	return config.DB.Where("project_id = ?", proj_id).Delete(&models.Resource{}).Error
}

func CheckResourceExistsById(resId string, exists *bool) error {
	return config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("id = ?", resId).Find(exists).Error
}

func CheckResourceExistsByKey(key string, projId string, exists *bool) error {
	return config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("key = ? AND project_id = ?",
		key, projId).Find(exists).Error
}
