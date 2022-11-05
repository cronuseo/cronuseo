package repositories

import (
	"fmt"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetResources(resources *[]models.Resource, proj_id string) {
	config.DB.Model(&models.Resource{}).Where("project_id = ?", proj_id).Find(&resources)
}

func GetResource(resource *models.Resource, res_id string) {
	config.DB.Where("id = ?", res_id).First(&resource)
}

func CreateResource(resource *models.Resource) {
	config.DB.Create(&resource)
}

func DeleteResource(resource *models.Resource, res_id string) {
	config.DB.Where("id = ?", res_id).Delete(&resource)
}

func UpdateResource(resource *models.Resource) {
	config.DB.Save(&resource)
}

func DeleteAllResources(proj_id string) {
	resources := []models.Resource{}
	GetResources(&resources, proj_id)
	for _, resource := range resources {
		res_id := resource.ID
		DeleteAllResourceActions(fmt.Sprint(res_id))
		DeleteAllResourceRoles(fmt.Sprint(res_id))
	}
	config.DB.Where("project_id = ?", proj_id).Delete(&models.Resource{})
}

func CheckResourceExistsById(resId string, exists bool) error {
	return config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("id = ?", resId).Find(&exists).Error
}

func CheckResourceExistsByKey(key string, proj_id string, exists bool) error {
	return config.DB.Model(&models.Resource{}).Select("count(*) > 0").Where("key = ? AND project_id = ?", key, proj_id).Find(&exists).Error
}
