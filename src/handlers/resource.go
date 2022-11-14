package handlers

import (
	"github.com/shashimalcse/Cronuseo/models"
	"github.com/shashimalcse/Cronuseo/repositories"
)

func GetResources(project_id string, resources *[]models.Resource) error {
	return repositories.GetResources(project_id, resources)
}

func GetResource(project_id string, id string, resource *models.Resource) error {
	return repositories.GetResource(project_id, id, resource)
}

func CreateResource(project_id string, resource *models.Resource) error {
	return repositories.CreateResource(project_id, resource)
}

func DeleteResource(project_id string, id string) error {
	return repositories.DeleteResource(project_id, id)
}

func UpdateResource(project_id string, id string, resource *models.Resource,
	reqResource *models.ResourceUpdateRequest) error {
	err := repositories.GetResource(project_id, id, resource)
	if err != nil {
		return err
	}
	resource.Name = reqResource.Name
	return repositories.UpdateResource(resource)
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
