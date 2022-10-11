package repositories

import (
	"errors"

	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetProjects(projects *[]models.Project, org_id string) {
	config.DB.Model(&models.Project{}).Where("organization_id = ?", org_id).Find(&projects)
}

func GetProject(project *models.Project, proj_id string) {
	config.DB.Where("id = ?", proj_id).First(&project)
}

func CreateProject(project *models.Project) {
	config.DB.Create(&project)
}

func DeleteProject(project *models.Project, proj_id string) {
	DeleteAllResources(proj_id)
	config.DB.Where("id = ?", proj_id).Delete(&project)
}

func UpdateProject(project *models.Project, reqProject *models.Project, proj_id string) {
	config.DB.Where("id = ?", proj_id).First(&project)
	project.Name = reqProject.Name
	project.Description = reqProject.Description
	config.DB.Save(&project)
}

func CheckProjectExistsById(proj_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("id = ?", proj_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("project not exists")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckProjectExistsByKey(key string, org_id string) (bool, error) {
	var exists bool
	err := config.DB.Model(&models.Project{}).Select("count(*) > 0").Where("key = ? AND organization_id = ?", key, org_id).Find(&exists).Error
	if err != nil {
		return false, errors.New("")
	}
	if exists {
		return true, nil
	} else {
		return false, nil
	}
}
