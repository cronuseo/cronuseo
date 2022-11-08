package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetProjects(projects *[]models.Project, orgId string) error {
	return config.DB.Model(&models.Project{}).Where("organization_id = ?", orgId).Find(&projects).Error
}

func GetProject(project *models.Project, projId string) error {
	return config.DB.Where("id = ?", projId).First(&project).Error
}

func CreateProject(project *models.Project) error {
	return config.DB.Create(&project).Error
}

func DeleteProject(project *models.Project, projId string) error {
	return config.DB.Where("id = ?", projId).Delete(&project).Error
}

func UpdateProject(project *models.Project) error {
	return config.DB.Save(&project).Error
}

func CheckProjectExistsById(projId string, exists *bool) error {
	return config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("id = ?", projId).Find(exists).Error
}

func CheckProjectExistsByKey(key string, orgId string, exists *bool) error {
	return config.DB.Model(&models.Project{}).Select("count(*) > 0").Where(
		"key = ? AND organization_id = ?", key, orgId).Find(exists).Error
}
