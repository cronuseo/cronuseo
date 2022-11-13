package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResources(project_id string, projects *[]models.Resource) error {
	return repositories.GetResources(project_id, projects)
}

func GetResource(project_id string, id string, project *models.Resource) error {
	return repositories.GetResource(project_id, id, project)
}

func CreateResource(project_id string, project *models.Resource) error {
	return repositories.CreateResource(project_id, project)
}

func DeleteResource(project_id string, id string) error {
	return repositories.DeleteResource(project_id, id)
}

func UpdateResource(project_id string, id string, project *models.Resource,
	reqResource *models.ResourceUpdateRequest) error {
	err := repositories.GetResource(project_id, id, project)
	if err != nil {
		return err
	}
	project.Name = reqResource.Name
	return repositories.UpdateResource(project)
}

func CheckResourceExistsById(id string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceExistsById(id, &exists)
	return exists, err
}

func CheckResourceExistsByKey(project_id string, key string) (bool, error) {
	var exists bool
	err := repositories.CheckResourceExistsByKey(project_id, key, &exists)
	return exists, err
}
