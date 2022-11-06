package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetProjects(projects *[]models.Project, orgId string) {
	config.DB.Model(&models.Project{}).Where("organization_id = ?", orgId).Find(&projects)
}

func GetProject(project *models.Project, projId string) {
	config.DB.Where("id = ?", projId).First(&project)
}

func CreateProject(project *models.Project) {
	config.DB.Create(&project)
}

func DeleteProject(project *models.Project, projId string) {
	config.DB.Where("id = ?", projId).Delete(&project)
}

func UpdateProject(project *models.Project) {
	config.DB.Save(&project)
}

func CheckProjectExistsById(projId string, exists *bool) error {
	return config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("id = ?", projId).Find(exists).Error
}

func CheckProjectExistsByKey(key string, orgId string, exists *bool) error {
	return config.DB.Model(&models.Project{}).Select("count(*) > 0").Where(
		"key = ? AND organization_id = ?", key, orgId).Find(exists).Error
}
